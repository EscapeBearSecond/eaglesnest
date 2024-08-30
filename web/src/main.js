import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/loading/style/css'
import 'element-plus/es/components/notification/style/css'
import 'element-plus/es/components/message-box/style/css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './style/element_visiable.scss'

import { createApp } from 'vue'
import './core/admin'
// 引入封装的router
import router from '@/router/index'
import '@/permission'
import run from '@/core/admin.js'
import auth from '@/directive/auth'
import { store } from '@/pinia'
import App from './App.vue'
import { initDom } from './utils/positionToCode'
import Pagination from '@/components/Pagination/index.vue';
import AdvancedTable from '@/components/AdvancedTable/index.vue'

initDom()

const app = createApp(App)
app.config.productionTip = false

app.component('Pagination', Pagination)
app.component('advance-table', AdvancedTable)


app
    .use(run)
    .use(store)
    .use(auth)
    .use(router)
    .mount('#app')
export default app
