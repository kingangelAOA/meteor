import { createRouter, createWebHashHistory } from 'vue-router'
import Layout from '@/layout'

/* Router Modules */
import chartsRouter from './modules/charts'
import useExample from './modules/useExample'

export const constantRoutes = [
  {
    path: '/dashboard',
    component: Layout,
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        //using el svg icon, the elSvgIcon first when at the same time using elSvgIcon and icon
        meta: { title: 'Dashboard', elSvgIcon: 'Fold' }
      }
    ]
  },
  {
    path: '/flow',
    component: Layout,
    beforeEnter: (to, from) => {
      const name = from.name
      if (name === 'Pipeline') {
        to.params.from = 'Pipeline'
      } else if (name === 'Task') { 
        to.params.from = 'Task'
      }
    },
    children: [
      {
        path: 'index',
        name: 'Flow',
        component: () => import('@/views/pipeline/flow/index.vue'),
        meta: { title: 'flow', icon: 'flow', noCache: true },
        hidden: true
      }
    ]
  },
  {
    path: '/monitor',
    component: Layout,
    beforeEnter: (to, from) => {
      const name = from.name
      if (name === 'Pipeline') {
        to.params.from = 'Pipeline'
      } else if (name === 'Task') { 
        to.params.from = 'Task'
      }
    },
    children: [
      {
        path: 'index',
        name: 'Monitor',
        component: () => import('@/views/pipeline/monitor/index.vue'),
        meta: { title: 'monitor', icon: 'monitor', noCache: true },
        hidden: true
      }
    ]
  },
  {
    path: '/task',
    component: Layout,
    children: [
      {
        path: 'index',
        name: 'Task',
        component: () => import('@/views/pipeline/task/index.vue'),
        meta: { title: 'task', icon: 'task', noCache: true },
        hidden: true
      }
    ]
  },
  {
    path: '/',
    redirect: '/pipeline/index',
    component: Layout,
    children: [
      {
        path: 'pipeline/index',
        name: 'Pipeline',
        component: () => import('@/views/pipeline/index.vue'),
        meta: { title: 'pipeline', icon: 'example', noCache: true }
      }
    ]
  },
  {
    path: '/node',
    component: Layout,
    children: [
      {
        path: 'index',
        name: 'Node',
        component: () => import('@/views/node/index.vue'),
        meta: { title: 'node', icon: 'example', noCache: true }
      }
    ]
  },
  {
    path: '/plugin',
    component: Layout,
    children: [
      {
        path: 'index',
        name: 'Plugin',
        component: () => import('@/views/plugin/index.vue'),
        meta: { title: 'plugin', icon: 'example', noCache: true }
      },
      {
        path: '/plugin/edit',
        hidden: true,
        name: 'PluginEdit',
        component: () => import('@/views/plugin/edit.vue'),
        meta: { title: 'pluginEdit', icon: 'example', noCache: true }
      }
    ]
  },
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login/Login.vue'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404.vue'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401.vue'),
    hidden: true
  }
  
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  // {
  //   path: '/permission',
  //   component: Layout,
  //   redirect: '/permission/page',
  //   alwaysShow: true, // will always show the root menu
  //   name: 'Permission',
  //   meta: {
  //     title: 'Permission',
  //     icon: 'lock',
  //     roles: ['admin', 'editor'] // you can set roles in root nav
  //   },
  //   children: [
  //     {
  //       path: 'roleIndex',
  //       component: () => import('@/views/permission'),
  //       name: 'Permission',
  //       meta: {
  //         title: 'Role Index'
  //         //roles: ['admin'] // or you can only set roles in sub nav
  //       }
  //     },
  //     {
  //       path: 'page',
  //       component: () => import('@/views/permission/page.vue'),
  //       name: 'PagePermission',
  //       meta: {
  //         title: 'Page Permission',
  //         roles: ['admin'] // or you can only set roles in sub nav
  //       }
  //     },
  //     {
  //       path: 'directive',
  //       component: () => import('@/views/permission/directive.vue'),
  //       name: 'DirectivePermission',
  //       meta: {
  //         title: 'Directive Permission'
  //         // if do not set roles, means: this page does not require permission
  //       }
  //     },
  //     {
  //       path: 'code-index',
  //       component: () => import('@/views/permission/CodePermission.vue'),
  //       name: 'CodePermission',
  //       meta: {
  //         title: 'Code Index'
  //       }
  //     },
  //     {
  //       path: 'code-page',
  //       component: () => import('@/views/permission/CodePage.vue'),
  //       name: 'CodePage',
  //       meta: {
  //         title: 'Code Page',
  //         code: 1
  //       }
  //     }
  //   ]
  // },
  // 404 page must be placed at the end !!!
  // using pathMatch install of "*" in vue-router 4.0
  { path: '/:pathMatch(.*)', redirect: '/404', hidden: true }
]

const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior: () => ({ top: 0 }),
  routes: constantRoutes
})

export default router
