<template>
    <div class="clear" style="clear: both">
        <div class="item">
            <div class="name">
                {{name}}
            </div>
            <img :src="head" alt="" class="head">
            <div v-if="img">
                <img :src="img" alt="" class="img">
            </div>
            <span v-if="msg">
                {{msg}}
            </span>
        </div>
        <mu-icon-button icon="delete" @click="del(id)" slot="right"/>
    </div>
</template>

<script type="text/ecmascript-6">
export default{
    props: ['name', 'img', 'msg', 'head','id'],
    methods: {
        del(open_id) {
            if (open_id != 0) {
                this.$store.state.msg_id = open_id
                var data = {
                    method: "delete_msg",
                    id: open_id,
                }
                this.$store.state.socket.send(JSON.stringify(data));
            }
        }
    }
}
</script>

<style lang="stylus" rel="stylesheet/stylus" scoped>
.clear
    margin-top :10px
    .item
        position: relative
        clear: both
        display: inline-block
        padding: 16px 40px 16px 20px
        margin: 10px 10px 20px 10px
        border-radius: 10px
        background-color: rgba(25, 147, 147, 0.2)
        animation: show-chat-odd 0.25s 1 ease-in
        -moz-animation: show-chat-odd 0.25s 1 ease-in
        -webkit-animation: show-chat-odd 0.25s 1 ease-in
        float: right
        margin-right: 80px
        color: #0AD5C1
        .img
            max-width: 200px
        span
            word-break:break-all
        .name
            position: absolute
            top: -20px
            width: 50px
            height: 20px
            right: -70px;
            text-align center
        .head
            position: absolute
            top: 0
            width: 50px
            height: 50px
            border-radius: 50px
            right: -70px;
        &:after
            position: absolute
            top: 15px
            content: ''
            width: 0
            height: 0
            border-top: 15px solid rgba(25, 147, 147, 0.2)
            border-right: 15px solid transparent
            right: -15px
@keyframes show-chat-odd {
    0% {
        margin-right: -480px;
    }

    100% {
        margin-right: 0;
    }
}
@-moz-keyframes show-chat-odd {
    0% {
        margin-right: -480px;
    }

    100% {
        margin-right: 0;
    }
}
@-webkit-keyframes show-chat-odd {
    0% {
        margin-right: -480px;
    }

    100% {
        margin-right: 0;
    }
}
</style>
