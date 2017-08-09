import Vue from 'vue'
import 'muse-components/styles/base.less' // 全局样式包含 normalize.css

import appBar from 'muse-components/appBar'
import avatar from 'muse-components/avatar'
import iconButton from 'muse-components/iconButton'
import { swiper, swiperSlide } from 'vue-awesome-swiper'
import icon from 'muse-components/icon'
import flatButton from 'muse-components/flatButton'
import raisedButton from 'muse-components/raisedButton'
import paper from 'muse-components/paper'
import chip from 'muse-components/chip'
import contentBlock from 'muse-components/contentBlock'
import dialog from 'muse-components/dialog'
import divider from 'muse-components/divider'
import badge from 'muse-components/badge'
import subHeader from 'muse-components/subHeader'
import { card, cardTitle, cardMedia, cardText, cardActions, cardHeader } from 'muse-components/card'
import App from './App.vue'
import VueRouter from "vue-router";
Vue.use(VueRouter);

import { bottomNav, bottomNavItem } from 'muse-components/bottomNav'
import { gridList, gridTile } from 'muse-components/gridList'
import { flexbox, flexboxItem } from 'muse-components/flexbox'
import { list, listItem } from 'muse-components/list'
import textField from 'muse-components/textField'

import 'muse-ui/dist/theme-light.css' // 使用 light 主题
import { tabs, tab } from 'muse-components/tabs'


// ..
Vue.component(appBar.name, appBar)
Vue.component(avatar.name, avatar)
Vue.component(flatButton.name, flatButton)
Vue.component(iconButton.name, iconButton)
Vue.component(raisedButton.name, raisedButton)
Vue.component(icon.name, icon)
Vue.component(bottomNav.name, bottomNav)
Vue.component(bottomNavItem.name, bottomNavItem)
Vue.component(paper.name, paper)
Vue.component(contentBlock.name, contentBlock)
Vue.component(card.name, card)
Vue.component(cardMedia.name, cardMedia)
Vue.component(cardTitle.name, cardTitle)
Vue.component(cardText.name, cardText)
Vue.component(cardHeader.name, cardHeader)
Vue.component(cardActions.name, cardActions)
Vue.component(chip.name, chip)
Vue.component(textField.name, textField)
Vue.component(badge.name, badge)
Vue.component(dialog.name, dialog)
Vue.component(divider.name, divider)
Vue.component(subHeader.name, subHeader)

Vue.component(gridList.name, gridList)
Vue.component(gridTile.name, gridTile)
Vue.component(flexbox.name, flexbox)
Vue.component(flexboxItem.name, flexboxItem)
Vue.component(list.name, list)
Vue.component(listItem.name, listItem)
Vue.component(tabs.name, tabs)
Vue.component(tab.name, tab)

Vue.component(swiper.name, swiper)
Vue.component(swiperSlide.name, swiperSlide)

// 创建一个路由器实例
// 并且配置路由规则
const router = new VueRouter({
    //mode: 'history',
    base: __dirname,
    routes: [{
        path: '/',
        component: (resolve) => {
            require.ensure([], () => resolve(require('./components/Login.vue')), 'login');
        }
    }, {
        path: '/contacts',
        component: (resolve) => {
            require.ensure([], () => resolve(require('./components/Contacts.vue')), 'contacts');
        }
    }, {
        path: '/index',
        component: (resolve) => {
            require.ensure([], () => resolve(require('./components/Login.vue')), 'index');
        }
    }]
})

new Vue({
    router: router,
    el: '#app',
    render: h => h(App)
})
