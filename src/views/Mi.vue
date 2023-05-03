<template>
    <v-navigation-drawer v-model="show_drawer" app>
        <sidebar :option="option" />
    </v-navigation-drawer>

    <v-app-bar app color="indigo" flat dark height="50px">
        <v-app-bar-nav-icon @click.stop="show_drawer = !show_drawer" />
        <v-toolbar-title>mi</v-toolbar-title>
        <v-spacer />
    </v-app-bar>

    <v-main>
        <v-container>
            <v-row class="boards_wrap">
                <v-col class="board_wrap" cols="auto" v-for="board_info in opened_board_infos"
                    :key="board_info.board_info_id">
                    <board :board_info="board_info" :task_infos="new Array<TaskInfo>()" /> <!-- これどうすんの？ -->
                    <!-- //TODO @諸々飛んできたときの処理 -->
                </v-col>
            </v-row>
            <v-row class="detail_task_row">
                <v-col class="detail_task_wrap" cols="auto">
                    <detail_task v-if="watching_task_info != null" :task_info="watching_task_info" />
                    <!-- //TODO @諸々飛んできたときの処理 -->
                </v-col>
            </v-row>
        </v-container>
    </v-main>

    <v-snackbar v-model="show_message_snackbar">
        <v-row>
            <v-col cols="auto">
                {{ message }}
            </v-col>
            <v-btn icon @click="show_message_snackbar = false">
                <v-icon>mdi-close</v-icon>
            </v-btn>
        </v-row>
    </v-snackbar>
</template>

<script setup lang="ts">
import { Ref, ref, watch, nextTick } from 'vue';
import sidebar from './sidebar/sidebar.vue';
import board from './board/board.vue';
import detail_task from './task/detail_task.vue';
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';
import BoardInfo from '@/api/data_struct/BoardInfo';
import TaskInfo from '@/api/data_struct/TaskInfo';

const show_drawer: Ref<boolean | null> = ref(null)
const show_message_snackbar: Ref<boolean> = ref(false)
const message: Ref<string> = ref("")
const option: Ref<ApplicationConfig> = ref(new ApplicationConfig())
const opened_board_infos: Ref<Array<BoardInfo>> = ref(new Array<BoardInfo>())
const watching_task_info: Ref<TaskInfo | null> = ref(null)
</script>

<style></style>