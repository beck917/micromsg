<template>
    <section class="main-content">
        <form action="" name="searchform">
        <mu-appbar title="联系人">

            <mu-text-field icon="search" name="search" class="appbar-search-field"  slot="right" hintText="请输入搜索内容"/>
            <mu-flat-button color="white" @click="add" label="添加" slot="right"/>

        </mu-appbar>
        </form>
        <mobile-tear-sheet>
            <mu-list>
                <mu-sub-header>最近聊天记录</mu-sub-header>
                <template v-for="item in this.$store.state.contacts">
                <mu-list-item  :title="item.cname" @click="open(item.cid)">
                    <mu-avatar :src="'/dist/images/a'+item.cid+'.jpg'" slot="leftAvatar"/>
                    <mu-badge :content="item.unread" circle secondary slot="right"/>
                    <mu-list-item slot="nested" title="删除" @click="del(item.cid)">
                        <mu-icon slot="right" value="delete"/>
                    </mu-list-item>
                </mu-list-item>
                </template>
            </mu-list>
        </mobile-tear-sheet>
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
                var name = document.searchform.search.value.trim()
                if (name != "") {
                    var data = {
                        method: "add",
                        add_name: name,
                    }
                    this.$store.state.socket.send(JSON.stringify(data));
                }
            },
            del(open_id) {
                if (open_id != 0) {
                    this.$store.state.open_id = open_id
                    var data = {
                        method: "delete",
                        delete_id: open_id,
                    }
                    this.$store.state.socket.send(JSON.stringify(data));
                }
            }
        }
    }
</script>
<style type="text/css">

</style>

