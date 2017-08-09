<template>
    <section class="main-content">
        <div class="center">
            <mu-text-field label="用户名" hintText="请输入用户名" labelFloat  icon="phone"/><br/>
            <mu-text-field label="密码" hintText="请输入密码" type="password" labelFloat icon="comment"/><br/>
            <mu-raised-button label="登录" v-on:click="login" class="demo-raised-button"/>
        </div>
    </section>
</template>
<script>
    export default {
    data() {
        return {}
    },
    methods: {
        login() {
            this.$router.push('contacts')
        }
    }
}
</script>
<style type="text/css">
    .center {
        padding: 40% 0;
        text-align:center;
    }
</style>

