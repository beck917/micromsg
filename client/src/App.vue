<template>
    <div id="app">
        <header class="header">
            <!--
            <mu-appbar :title="current" class="common-top-bar">
                <mu-icon-button icon="menu" slot="left"/>
            </mu-appbar>
            -->
        </header>
        <router-view class="view"></router-view>
        <footer class="footer">
            <!--
            <mu-paper>
                <mu-bottom-nav :value="bottomNav" @change="handleChange">
                    <mu-bottom-nav-item to="index" value="recents" title="主页" icon="home" />
                    <mu-bottom-nav-item value="search" title="发现" icon="search" />
                    <mu-bottom-nav-item value="camera" title="照相" icon="camera" />
                    <mu-bottom-nav-item value="favorites" title="收藏" icon="favorite" />
                    <mu-bottom-nav-item to="person" value="nearby" title="我的" icon="account_circle" />
                </mu-bottom-nav>
            </mu-paper>
            -->
        </footer>
    </div>
</template>
<script>
export default {
    data() {
            return {
                bottomNav: 'recents',
                current: '首页',
                user: {
                    username:'咸蛋超人',
                    avatar: 'https://www.muse-ui.org/images/uicon.jpg',
                }
            }
        },
        methods: {
            handleChange(val) {
                this.bottomNav = val
            }
        }
}
</script>
<style lang="scss">
html {
    height: 100%;
}

body {
    height: 100%;
}

#app {
    display: flex;
    flex-direction: column;
    height: 100%;
}

.header {
    /* 我们希望 header 采用固定的高度，只占用必须的空间 */
    /* 0 flex-grow, 0 flex-shrink, auto flex-basis */
    flex: 0 0 auto;
}
.header .mu-text-field{
    margin-bottom: 0;
}
.header .user-avatar{
    line-height: 56px;
    vertical-align: middle;
}
.main-content {
    /* 将 flex-grow 设置为1，该元素会占用所有的可使用空间
      而其他元素该属性值为0，因此不会得到多余的空间*/
    /* 1 flex-grow, 0 flex-shrink, auto flex-basis */
    flex: 1 1 auto;
    overflow-y: scroll;
}

.footer {
    display: block;
    /* 和 header 一样，footer 也采用固定高度*/
    /* 0 flex-grow, 0 flex-shrink, auto flex-basis */
    flex: 0 0 auto;
}


/* fallback */

@font-face {
    font-family: 'Material Icons';
    font-style: normal;
    font-weight: 400;
    src: local('Material Icons'), local('MaterialIcons-Regular'), url(https://fonts.gstatic.com/s/materialicons/v22/2fcrYFNaTjcS6g4U3t-Y5ZjZjT5FdEJ140U2DJYC3mY.woff2) format('woff2');
}

.material-icons {
    font-family: 'Material Icons';
    font-weight: normal;
    font-style: normal;
    font-size: 24px;
    line-height: 1;
    letter-spacing: normal;
    text-transform: none;
    display: inline-block;
    white-space: nowrap;
    word-wrap: normal;
    direction: ltr;
    -webkit-font-feature-settings: 'liga';
    -webkit-font-smoothing: antialiased;
}
    .common-top-bar{
        /*text-align: center;*/
    }
</style>

<style lang="stylus" rel="stylesheet/stylus">
    #app
        width: 100%
        min-height: 100%
        .bottom
            position: fixed
            bottom: 0
            width: 100%


</style>
