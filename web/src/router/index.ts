import { createRouter, createWebHistory } from 'vue-router'

import Login from '../pages/login/Login.vue'

// ==================== cluster ====================
import Node from '../pages/cluster/Node.vue'
import Slot from '../pages/cluster/Slot.vue'
import Channel from '../pages/cluster/Channel.vue'
import Config from '../pages/cluster/Config.vue'
import Log  from '../pages/cluster/Log.vue'
// import NodeDetail from '../pages/cluster/NodeDetail.vue'

// ==================== data ====================
import User from '../pages/data/User.vue'
import Device from '../pages/data/Device.vue'
import Connection from '../pages/data/Connection.vue' 
import Message from '../pages/data/Message.vue'
import ChannelForData from '../pages/data/Channel.vue'
import Conversation from '../pages/data/Conversation.vue'

// ==================== monitor ====================

import MonitorApp from '../pages/monitor/App.vue'
import MonitorCluster from '../pages/monitor/Cluster.vue'
import MonitorSystem from '../pages/monitor/System.vue'
import MonitorDB from '../pages/monitor/DB.vue'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // ==================== cluster ====================
    {
      path: '/',
      name: 'home',
      component: Node
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/cluster/nodes',
      name: 'node',
      component: Node
    },
    // {
    //   path: '/cluster/node/detail',
    //   name: 'nodeDetail',
    //   component: NodeDetail
    // },
    {
      path: '/cluster/slots',
      name: 'slot',
      component: Slot
    },
    {
      path: '/cluster/channels',
      name: 'channel',
      component: Channel
    },
    {
      path: '/cluster/config',
      name: 'config',
      component: Config
    },
    {
      path: '/cluster/log',
      name: 'log',
      component: Log
    },

    // ==================== data ====================
    {
      path: '/data/connection',
      name: 'dataConnection',
      component: Connection
    },
    {
      path: '/data/user',
      name: 'dataUser',
      component: User
    },
    {
      path: '/data/device',
      name: 'dataDevice',
      component: Device
    },
    {
      path: '/data/message',
      name: 'dataMessage',
      component: Message
    },
    {
      path: '/data/channel',
      name: 'dataChannel',
      component: ChannelForData
    },
    {
      path: '/data/conversation',
      name: 'dataConversation',
      component: Conversation
    },
      // ==================== monitor ====================
    {
      path: '/monitor/app',
      name: 'monitorApp',
      component: MonitorApp
    },
    {
      path: '/monitor/cluster',
      name: 'monitorCluster',
      component: MonitorCluster
    },
    {
      path: '/monitor/system',
      name: 'monitorSystem',
      component: MonitorSystem
    },
    {
      path: '/monitor/db',
      name: 'monitorDB',
      component: MonitorDB
    },
  ]
})

export default router
