<template>
    <section class="main-content">
        <mobile-tear-sheet>
            <mu-list>
                <mu-sub-header>最近聊天记录</mu-sub-header>
                <mu-list-item title="Mike Li" @click="open">
                    <mu-avatar src="/images/a1.jpg" slot="leftAvatar"/>
                    <mu-icon value="chat_bubble" slot="right"/>
                </mu-list-item>
                <mu-list-item title="Maco Mai">
                    <mu-avatar src="../assets/avatar/a2.jpg" slot="leftAvatar"/>
                    <mu-icon value="chat_bubble" slot="right"/>
                </mu-list-item>
                <mu-list-item title="Alex Qin">
                    <mu-avatar src="/avatar/a3.jpg" slot="leftAvatar"/>
                    <mu-badge content="12" secondary slot="right"/>
                </mu-list-item>
                <mu-list-item title="Allen Qun">
                    <mu-avatar src="/avatar/a4.jpg" slot="leftAvatar"/>
                    <mu-badge content="12" circle secondary slot="right"/>
                </mu-list-item>
                <mu-list-item title="Myron Liu">
                    <mu-avatar src="/avatar/a5.jpg" slot="leftAvatar"/>
                    <mu-badge content="12" circle secondary slot="right">
                        <mu-icon-button icon="chat_bubble"/>
                    </mu-badge>
                </mu-list-item>
            </mu-list>
            <mu-divider/>
            <mu-list>
                <mu-sub-header>历史聊天记录</mu-sub-header>
                <mu-list-item title="Gaia Zhou">
                    <mu-avatar src="/images/avatar5.jpg" slot="leftAvatar"/>
                    <mu-icon value="chat_bubble" slot="right"/>
                </mu-list-item>
                <mu-list-item title="Miller Wang">
                    <mu-avatar src="/images/avatar6.jpg" slot="leftAvatar"/>
                    <mu-icon value="chat_bubble" slot="right"/>
                </mu-list-item>
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
            open() {
                this.dialog = true
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

