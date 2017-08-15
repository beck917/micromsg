<template>
        <div class="container">
            <div class="title">
                <mu-appbar title="聊天">
                    <mu-icon-button icon="chevron_left" slot="left" @click="close"/>
                </mu-appbar>

            </div>
            <div style="height:50px"></div>
            <div class="chat">
                <div v-for="obj in this.$store.state.msg_list">
                    <othermsg v-if="obj.send_uid != uid" :head="'/dist/images/a'+obj.send_uid+'.jpg'" :msg="obj.msg"
                              ></othermsg>
                    <mymsg v-if="obj.send_uid == uid" :id="obj.id" :head="'/dist/images/a'+obj.send_uid+'.jpg'" :msg="obj.msg"
                    ></mymsg>
                </div>
                <div style="height:200px"></div>
            </div>
            <div class="bottom">
                <div class="chat">
                    <div class="input" @keyup.enter="send">
                        <input type="text" id="message">
                    </div>
                    <mu-raised-button label="发送" class="demo-raised-button" primary  @click="send"/>
                </div>
            </div>
        </div>
</template>

<script>
    import Mymsg from './Mymsg.vue'
    import Othermsg from './Othermsg.vue'
    export default {
        data() {
            return {
                uid :  this.$store.state.uid,
            }
        },
        methods: {
            send() {
                // 判断发送信息是否为空
                if (document.getElementById('message').value !== '') {
                    var data = {
                        method: "send",
                        msg: document.getElementById('message').value,
                        send_id: this.$store.state.uid,
                        recv_id: this.$store.state.open_id,
                    }
                    this.$store.state.socket.send(JSON.stringify(data));

                    var msg_data = {
                        msg: document.getElementById('message').value,
                        send_uid: this.$store.state.uid,
                        recv_uid: this.$store.state.open_id,
                    }
                    this.$store.state.msg_list.push(msg_data)
                    window.scrollTo(0, 900000)
                    document.getElementById('message').value = ''
                } else {

                }
            },
            close() {
                this.$store.state.open_id = 0
                this.$router.back(-1)
            },
        },
        components: {
            Mymsg,
            Othermsg
        },
        mounted(){
            //  挂载结束状态
            window.scrollTo(0, 900000)
        }
    }
</script>

<style lang="stylus" rel="stylesheet/stylus" scoped>
    &.fade-enter-active, &.fade-leave-active
        transition: all 0.2s linear
        transform translate3d(0, 0, 0)

    &.fade-enter, &.fade-leave-active
        opacity: 1
        transform translate3d(100%, 0, 0)

    .container
        position: absolute
        left: 0
        top: 0
        width: 100%
        min-height: 100%
        background: #ffffff
        .title
            position: fixed
            height: 30px
            top: 0
            width: 100%
            z-index: 1
            .center
                -webkit-box-flex: 1
                -webkit-flex: 1
                -ms-flex: 1
                flex: 1
                padding-left: 8px
                padding-right: 8px
                white-space: nowrap
                text-overflow: ellipsis
                overflow: hidden
                font-size: 20px
                font-weight: 400
                line-height: 56px
                text-align: center
        .chat
        .all-chat
            .online
                display: inline-block
                margin: 5px
                img
                    width: 40px
                    height: 40px
                    border-radius: 100%
        .bottom
            position: fixed
            height: 50px
            bottom: 0
            background: #eeeff3
            .chat
                width: 100%
                display: flex
                .input
                    flex: 1
                    background: #eeeff3
                    padding: 4px
                    input
                        width: 100%
                        height: 42px
                        box-sizing: border-box
                        border: 1px solid #8c8c96
                        color: #333333
                        font-size: 18px
                        padding-left: 5px
                    .mu-text-field
                        width: 100%
                .demo-raised-button
                    flex-basis: 88px
                    margin-top: 4px
                    height: 40px
                    background: #eeeff3
                    color: #8c8c96
            .functions
                width: 100%
                .fun-li
                    width: 40px
                    height: 30px
                    display: inline-block
                .fun-li:nth-child(1)
                    background-image: url(../assets/images.png)
                    background-repeat: no-repeat
                    background-size: 25px 25px
                    background-position: center center

</style>
