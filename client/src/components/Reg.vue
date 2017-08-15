<template>
    <section class="login">
        <div class="header">
            <mu-appbar title="注册">
                <mu-flat-button label="注册" slot="default"/>
            </mu-appbar>
        </div>
        <div class="content">
            <form action="" name="loginform">
                <mu-text-field label="用户名" hintText="请输入用户名" name="username" labelFloat  icon="phone"/><br/>
                <mu-text-field label="密码" hintText="请输入密码" name="password" type="password" labelFloat icon="comment"/><br/>
                <mu-raised-button label="注册" v-on:click="reg" class="demo-raised-button"/>
            </form>
            <br />
            <div @click="loginu" class="tip-user">
                我已有帐号
            </div>
        </div>
    </section>
</template>
<script>
    export default {
        data() {
            return {}
        },
        methods: {
            reg() {
                var name = document.loginform.username.value.trim()
                var password = document.loginform.password.value.trim()
                if (name !== '' && password !== '') {
                    var data = {
                        method: "register",
                        username: name,
                        password: password
                    }

                    this.$store.state.socket.send(JSON.stringify(data));
                }
            },
            loginu() {
                this.$router.push("login")
            }
        }
    }
</script>
<style lang="stylus" rel="stylesheet/stylus">
    .login
        position: absolute
        left: 0
        right: 0
        top: 0
        bottom: 0
        background: #fff
        .mu-appbar
            text-align: center
            .mu-flat-button-label
                font-size: 20px
        .content
            width: 80%
            margin: 20px auto
            .mu-text-field
                width: 100%
            .mu-raised-button
                min-width: 80px
                display: block
                width: 100%
                margin: 0 auto
                transition: all 0.375s
                &.loading
                    width: 80px
                    height: 80px
                    border-radius: 50%
</style>

