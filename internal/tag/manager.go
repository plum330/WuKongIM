package tag

import (
	"hash/fnv"
	"time"

	"github.com/WuKongIM/WuKongIM/internal/errors"
	"github.com/WuKongIM/WuKongIM/internal/types"

	"github.com/WuKongIM/WuKongIM/internal/service"
	"github.com/WuKongIM/WuKongIM/pkg/wkutil"
)

type TagMgr struct {
	bluckets []*blucket
}

func NewTagMgr(blucketCount int) *TagMgr {
	tg := &TagMgr{}
	tg.bluckets = make([]*blucket, blucketCount)
	for i := 0; i < blucketCount; i++ {
		tg.bluckets[i] = newBlucket(i, time.Minute*20)
	}
	return tg
}

func (t *TagMgr) Start() error {
	var err error
	for _, b := range t.bluckets {
		err = b.start()
		if err != nil {
			return err
		}
	}
	return nil
}
func (t *TagMgr) Stop() {
	for _, b := range t.bluckets {
		b.stop()
	}
}

func (t *TagMgr) MakeTag(uids []string) (*types.Tag, error) {
	tagKey := wkutil.GenUUID()
	return t.MakeTagWithTagKey(tagKey, uids)
}

func (t *TagMgr) MakeTagWithTagKey(tagKey string, uids []string) (*types.Tag, error) {

	tag := &types.Tag{
		Key:         tagKey,
		LastGetTime: time.Now(),
	}

	nodes, err := t.calcUsersInNode(uids)
	if err != nil {
		return nil, err
	}
	tag.Nodes = nodes

	return tag, nil
}

func (t *TagMgr) AddUsers(tagKey string, uids []string) error {
	tag := t.getTag(tagKey)
	if tag == nil {
		return errors.TagNotExist(tagKey)
	}
	// 去除已经存在的用户
	t.removeExistUidsInTag(tag, uids)

	// 计算用户所在的节点
	nodes, err := t.calcUsersInNode(uids)
	if err != nil {
		return err
	}
	// 合并节点
	t.mergeNodes(tag, nodes)

	return nil
}

func (t *TagMgr) removeExistUidsInTag(tag *types.Tag, uids []string) {
	for _, node := range tag.Nodes {
		for _, uid := range uids {
			for i, nodeUid := range node.Uids {
				if nodeUid == uid {
					node.Uids = append(node.Uids[:i], node.Uids[i+1:]...)
					break
				}
			}
		}
	}
}

func (t *TagMgr) RemoveUsers(tagKey string, uids []string) error {
	tag := t.getTag(tagKey)
	if tag == nil {
		return errors.TagNotExist(tagKey)
	}

	for _, uid := range uids {
		slotId := service.Cluster.GetSlotId(uid)
		leaderId, err := service.Cluster.SlotLeaderId(slotId)
		if err != nil {
			return err
		}
		if leaderId == 0 {
			return errors.TagSlotLeaderIsZero
		}
		for _, node := range tag.Nodes {
			if node.LeaderId == leaderId {
				for i, nodeUid := range node.Uids {
					if nodeUid == uid {
						node.Uids = append(node.Uids[:i], node.Uids[i+1:]...)
						break
					}
				}
				break
			}
		}
	}
	return nil
}

func (t *TagMgr) RemoveTag(tagKey string) {
	t.removeTag(tagKey)
}

func (t *TagMgr) GetUsers(tagKey string) []string {
	tag := t.getTag(tagKey)
	if tag == nil {
		return nil
	}
	var uids []string
	for _, node := range tag.Nodes {
		uids = append(uids, node.Uids...)
	}
	return uids
}

func (t *TagMgr) Get(tagKey string) *types.Tag {
	tag := t.getTag(tagKey)
	tag.LastGetTime = time.Now()
	return tag
}

func (t *TagMgr) Exist(tagKey string) bool {
	return t.getTag(tagKey) != nil
}

func (t *TagMgr) RenameTag(oldTagKey, newTagKey string) error {
	tag := t.getTag(oldTagKey)
	if tag == nil {
		return errors.TagNotExist(oldTagKey)
	}
	tag.Key = newTagKey
	tag.LastGetTime = time.Now()
	t.setTag(tag)
	t.removeTag(oldTagKey)
	return nil
}

func (t *TagMgr) SetChannelTag(fakeChannelId string, channelType uint8, tagKey string) {
	blucket := t.getBlucketByChannel(fakeChannelId, channelType)
	blucket.setChannelTag(fakeChannelId, channelType, tagKey)
}

func (t *TagMgr) GetChannelTag(fakeChannelId string, channelType uint8) string {
	blucket := t.getBlucketByChannel(fakeChannelId, channelType)
	return blucket.getChannelTag(fakeChannelId, channelType)
}

func (t *TagMgr) getBlucketByTagKey(tagKey string) *blucket {
	h := fnv.New32a()
	h.Write([]byte(tagKey))
	i := h.Sum32() % uint32(len(t.bluckets))
	return t.bluckets[i]
}

func (t *TagMgr) getBlucketByChannel(channelId string, channelType uint8) *blucket {
	h := fnv.New32a()
	h.Write([]byte(wkutil.ChannelToKey(channelId, channelType)))
	i := h.Sum32() % uint32(len(t.bluckets))
	return t.bluckets[i]
}

func (t *TagMgr) mergeNodes(tag *types.Tag, nodes []*types.Node) {
	for _, node := range nodes {
		exist := false
		for _, tagNode := range tag.Nodes {
			if tagNode.LeaderId == node.LeaderId {
				exist = true

				// 合并用户
				existUser := false
				for _, uid := range node.Uids {
					for _, tagUid := range tagNode.Uids {
						if tagUid == uid {
							existUser = true
							break
						}
					}
					if !existUser {
						tagNode.Uids = append(tagNode.Uids, uid)
					}
				}
				// 合并slot
				for _, slotId := range node.SlotIds {
					existSlot := false
					for _, tagSlotId := range tagNode.SlotIds {
						if tagSlotId == slotId {
							existSlot = true
							break
						}
					}
					if !existSlot {
						tagNode.SlotIds = append(tagNode.SlotIds, slotId)
					}
				}

				break
			}
		}
		if !exist {
			tag.Nodes = append(tag.Nodes, node)
		}
	}
}

func (t *TagMgr) calcUsersInNode(uids []string) ([]*types.Node, error) {

	var nodeMap = make(map[uint64]*types.Node)
	for _, uid := range uids {
		slotId := service.Cluster.GetSlotId(uid)
		leaderId, err := service.Cluster.SlotLeaderId(slotId)
		if err != nil {
			return nil, err
		}
		if leaderId == 0 {
			return nil, errors.TagSlotLeaderIsZero
		}
		node := nodeMap[leaderId]
		if node == nil {
			node = &types.Node{
				LeaderId: leaderId,
			}
			nodeMap[leaderId] = node
		}
		node.Uids = append(node.Uids, uid)
		existSlot := false
		for _, slot := range node.SlotIds {
			if slot == slotId {
				existSlot = true
				break
			}
		}
		if !existSlot {
			node.SlotIds = append(node.SlotIds, slotId)
		}

	}
	nodes := make([]*types.Node, 0, len(nodeMap))
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}
	return nodes, nil

}

func (t *TagMgr) setTag(tag *types.Tag) {
	blucket := t.getBlucketByTagKey(tag.Key)
	blucket.setTag(tag)
}

func (t *TagMgr) getTag(tagKey string) *types.Tag {
	blucket := t.getBlucketByTagKey(tagKey)
	return blucket.getTag(tagKey)
}

func (t *TagMgr) removeTag(tagKey string) {
	blucket := t.getBlucketByTagKey(tagKey)
	blucket.removeTag(tagKey)
}
