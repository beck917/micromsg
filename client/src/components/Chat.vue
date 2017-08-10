<template>
    <transition name="fade">
        <div class="container" v-show="getchattoggle">
            <div class="title">
                <mu-appbar title="Title">
                    <mu-icon-button icon="chevron_left" slot="left" @click="closechat"/>
                    <div class="center">
                        聊天({{Object.keys(getusers).length}})
                    </div>
                    <mu-icon-button icon="expand_more" slot="right" @click="shownotice"/>
                </mu-appbar>
            </div>
            <div class="all-chat">
                <div style="height:70px"></div>
                <div>在线人员</div>
                <div v-for="obj in getusers" class="online">
                    <img :src="obj.src" alt="">
                </div>
            </div>
            <div class="chat">
                <div v-for="obj in getmesshistoryinfos">
                    <othermsg v-if="obj.username!=getusername" :name="obj.username" :head="obj.src" :msg="obj.msg"
                              :img="obj.img"></othermsg>
                    <mymsg v-if="obj.username==getusername" :name="obj.username" :head="obj.src" :msg="obj.msg"
                           :img="obj.img"></mymsg>
                </div>
                <div v-for="obj in getinfos">
                    <othermsg v-if="obj.username!=getusername" :name="obj.username" :head="obj.src" :msg="obj.msg"
                              :img="obj.img"></othermsg>
                    <mymsg v-if="obj.username==getusername" :name="obj.username" :head="obj.src" :msg="obj.msg"
                           :img="obj.img"></mymsg>
                </div>
                <div style="height:120px"></div>
            </div>
            <div class="bottom">
                <div class="chat">
                    <div class="input" @keyup.enter="submess">
                        <input type="text" id="message">
                    </div>
                    <mu-raised-button label="发送" class="demo-raised-button" primary  @click="submess"/>
                </div>
                <div class="functions">
                    <div class="fun-li" @click="imgupload"></div>
                </div>
                <input id="inputFile" name='inputFile' type='file' multiple='mutiple' accept="image/*;capture=camera"
                       style="display: none" @change="fileup">
            </div>
        </div>
    </transition>
</template>