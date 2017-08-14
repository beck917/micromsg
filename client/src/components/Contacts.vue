<template>
    <section class="main-content">
        <mobile-tear-sheet>
            <mu-list>
                <mu-sub-header>最近聊天记录</mu-sub-header>
                <template v-for="item in this.$store.state.contacts">
                <mu-list-item :title="item.cname" @click="open(item.cid)">
                    <mu-avatar :src="'/dist/images/a'+item.cid+'.jpg'" slot="leftAvatar"/>
                    <mu-badge :content="item.unread" circle secondary slot="right"/>
                </mu-list-item>
                </template>
            </mu-list>
        </mobile-tear-sheet>
        <mu-dialog :open="dialog" @close="close" title="聊天" scrollable>
            <mu-list>
                <template v-for="item in list">
                    <mu-list-item :title="item">
                        <mu-avatar src="/images/avatar5.jpg" slot="leftAvatar"/>
                    </mu-list-item>
                    <mu-divider/>
                </template>
            </mu-list>
            <mu-text-field multiLine :rows="1" :rowsMax="6" hintText="请输入聊天内容" slot="actions"/><mu-flat-button @click="add" icon="add" label="发送" primary slot="actions"/><br/>
            <mu-flat-button primary label="关闭" @click="close" slot="actions"/>
        </mu-dialog>
    </section>
</template>
<script>
    export default {
        data() {
            const list = []
            for (let i = 0; i < 1; i++) {
                list.push('聊天记录聊天记录聊天记录聊天记录聊' + i)
            }
            return {
                list,
                dialog: false
            }
        },
        methods: {
            open(open_id) {
                this.dialog = false
                this.$store.state.open_id = open_id

                if (open_id != 0) {
                    var data = {
                        method: "open",
                        open_id: open_id,
                    }

                    this.$store.state.socket.send(JSON.stringify(data));
                }
            },
            close() {
                this.dialog = false
            },
            add() {
                this.list.push('add聊天记录聊天记录聊天记录聊天记录聊')
            }
        }
    }
</script>
<style type="text/css">

</style>

