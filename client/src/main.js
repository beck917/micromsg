import Vue from 'vue'
import Vuex from 'vuex'
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
import popup from 'muse-components/popup'
//import iconMenu from 'muse-components/iconMenu'
//import {menu, menuItem} from 'muse-components/menu'
import { card, cardTitle, cardMedia, cardText, cardActions, cardHeader } from 'muse-components/card'
import App from './App.vue'
import VueRouter from "vue-router";
Vue.use(VueRouter);
Vue.use(Vuex)

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
Vue.component(popup.name, popup)
//Vue.component(iconMenu.name, iconMenu)
//Vue.component(menuItem.name, menuItem)
//Vue.component(menu.name, menu)

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
    }, {
        path: '/chat',
        component: (resolve) => {
            require.ensure([], () => resolve(require('./components/Chat.vue')), 'chat');
        }
    }, {
        path: '/reg',
        component: (resolve) => {
            require.ensure([], () => resolve(require('./components/Reg.vue')), 'reg');
        }
    }, {
        path: '/login',
        component: (resolve) => {
            require.ensure([], () => resolve(require('./components/Login.vue')), 'login');
        }
    }]
})

var socket = {};
connect(); // 自动连接

var uid = 0
var open_id = 0

function connect() {
    // 连接OR断开socket
    if (!socket.readyState || socket.readyState != 1) {
        socket = new WebSocket("ws://" + location.hostname +":9127");
    } else {
        socket.close();
    }

    socket.onopen = function() {
        // Web Socket 已连接上，使用 send() 方法发送数据
        console.log("Web Socket 已连接上");
    };

    socket.onmessage = function(evt) {
        var res = JSON.parse(evt.data);
        store.state.poptop = true
        store.state.popmsg = res.msg
        setTimeout(() => {
            store.state.poptop = false
        }, 1500)
        // console.log(res);
        if (res.replymethod == 'login' && res.result == 1) {
            console.log("login success", res);
            store.state.contacts = res.data.contacts
            store.state.uid = res.data.uid
            store.state.contacts.sort(compare("unread"))
            router.push('contacts')
        }

        if(res.result == 1) {
            if (res.replymethod == 'send') {
                console.log("send", res);
            } else if (res.replymethod == 'register') {
                console.log("reg", res);
                router.push("login");
            } else if (res.replymethod == 'open') {
                console.log("open", res);
                for (var i = 0; i < store.state.contacts.length; i++) {
                    if (store.state.contacts[i].cid == store.state.open_id) {
                        store.state.contacts[i].unread = 0;
                    }
                }

                if (!res.data.msg_list) {
                    res.data.msg_list = []
                }

                store.state.msg_list = rev_arr(res.data.msg_list)
                router.push('chat')
            } else if (res.replymethod == 'add') {
                console.log("add", res);
                if (!store.state.contacts) {
                    store.state.contacts = [res.data.contact]
                } else {
                    store.state.contacts.push(res.data.contact);
                }
            } else if (res.replymethod == 'delete') {
                console.log("delete", res);
                for (var i = 0; i < store.state.contacts.length; i++) {
                    if (store.state.contacts[i].cid == store.state.open_id) {
                        store.state.contacts.splice(i, 1)
                    }
                }
                store.state.open_id = 0
            } else if (res.replymethod == 'delete_msg') {
                console.log("delete_msg", res);
                for (var i = 0; i < store.state.msg_list.length; i++) {
                    if (store.state.msg_list[i].id == store.state.msg_id) {
                        store.state.msg_list.splice(i, 1)
                    }
                }
            } else if (res.replymethod == 'pushmsg') {
                console.log("pushmsg", res);

                if (res.data.send_id == store.state.open_id) {
                    var msg_data = {
                        msg: res.data.msg,
                        send_uid: res.data.send_id,
                        recv_uid: res.data.recv_id,
                    }
                    store.state.msg_list.push(msg_data)
                    window.scrollTo(0, 900000)
                }
                //刷新联系人列表
                for (var i = 0; i < store.state.contacts.length; i++) {
                    if (store.state.contacts[i].cid == res.data.send_id) {
                        store.state.contacts[i].unread += 1;
                    }
                }
                if (res.data.contact) {
                    store.state.contacts.push(res.data.contact)
                }
                store.state.contacts.sort(compare("unread"))
            }
        }
    };

    socket.onclose = function() {
        console.error("Web Socket 已经断开");
    }
}

function rev_arr(arr) {
    var newarr = []
    for(var i=arr.length-1;i>=0;i--){
        newarr.push(arr[i])
    }
    return newarr
}

var compare = function (prop) {
    return function (obj1, obj2) {
        var val1 = obj1[prop];
        var val2 = obj2[prop];
        if (val1 < val2) {
            return 1;
        } else if (val1 > val2) {
            return -1;
        } else {
            return 0;
        }
    }
}

const store = new Vuex.Store({
    state: {
        // 存放用户
        socket: socket,
        uid:uid,
        open_id:open_id,
        contacts: [],
        msg_list:[],
        poptop:false,
        popmsg:"test",
        msg_id:0,
    }
})

new Vue({
    router: router,
    el: '#app',
    store,
    render: h => h(App)
})
