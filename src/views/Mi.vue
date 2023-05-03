<template>
    <v-navigation-drawer v-model="show_drawer" app>
        <sidebar :option="option" ref="sidebar_ref" @errors="write_messages"
            @updated_check_condition="updated_check_condition" @updated_sort_type="updated_sort_type"
            @updated_search_word="updated_search_word" @updated_boards_by_user="updated_boards_by_user"
            @clicked_board="clicked_board" @updated_tags_by_user="updated_tags_by_user" @updated_tags="updated_tags" />
    </v-navigation-drawer>

    <v-app-bar app color="indigo" flat dark height="50px">
        <v-app-bar-nav-icon @click.stop="show_drawer = !show_drawer" />
        <v-toolbar-title>mi</v-toolbar-title>
        <v-spacer />
    </v-app-bar>

    <v-main>
        <v-container>
            <v-row class="boards_wrap">
                <v-col class="board_wrap" cols="auto" v-for="board_name in opened_board_names" :key="board_name"
                    @copied_task_id="copied_task_id" @added_tag="added_tag" @added_text="added_text"
                    @updated_task="updated_task" @deleted_task="deleted_task">
                    <board :board_name="board_name" :task_infos="task_infos_map[board_name]" />
                </v-col>
            </v-row>
            <v-row class="detail_task_row">
                <v-col class="detail_task_wrap" cols="auto">
                    <detail_task v-if="watching_task_info != null" :task_info="watching_task_info"
                        @copied_task_id="copied_task_id" @added_tag="added_tag" @added_text="added_text"
                        @updated_task="updated_task" @deleted_task="deleted_task" @deleted_tag="deleted_tag"
                        @deleted_text="deleted_text" />
                </v-col>
            </v-row>
        </v-container>
    </v-main>

    <v-fab-transition>
        <v-btn fab dark color="indigo" class="mr-1 mb-10" absolute bottom right @click="show_add_task_dialog">
            <v-icon>mdi-plus</v-icon>
        </v-btn>
    </v-fab-transition>

    <add_task_dialog @errors="write_messages" @added_task="added_task" ref="add_task_dialog_ref" />

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
import add_task_dialog from './dialog/add_task_dialog.vue';
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';
import BoardInfo from '@/api/data_struct/BoardInfo';
import TaskInfo from '@/api/data_struct/TaskInfo';
import CheckState from '@/api/data_struct/CheckState';
import SortType from '@/api/data_struct/SortType';
import MiServerAPI from '@/api/MiServerAPI';
import GetApplicationConfigRequest from '@/api/GetApplicationConfigRequest';
import TaskSearchQuery from '@/api/data_struct/TaskSearchQuery';
import GetTasksFromBoardRequest from '@/api/GetTasksFromBoardRequest';

const show_drawer: Ref<boolean | null> = ref(null)
const show_message_snackbar: Ref<boolean> = ref(false)
const message: Ref<string> = ref("")
const option: Ref<ApplicationConfig> = ref(new ApplicationConfig())
const opened_board_names: Ref<Array<string>> = ref(new Array<string>())
const watching_task_info: Ref<TaskInfo | null> = ref(null)
const watching_board_name: Ref<string | null> = ref(null)
const sidebar_ref = ref<InstanceType<typeof sidebar> | null>(null);
const add_task_dialog_ref = ref<InstanceType<typeof add_task_dialog> | null>(null);
const query_map:Ref<any> = ref({})
const task_infos_map:Ref<any> = ref({})

update_option()

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))
function update_option(): void {
    const api = new MiServerAPI()
    const request = new GetApplicationConfigRequest()
    api.get_application_config(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                write_messages(res.errors)
                return
            }
            option.value = res.application_config
        })
}
function write_message(message_: string) {
    message.value = message_
    show_message_snackbar.value = true
}
function open_board(board_name: string) {
    const api = new MiServerAPI()
    opened_board_names.value.push(board_name)
    const query = sidebar_ref.value?.construct_task_search_query()
    query!.board = board_name
    query_map.value[board_name] = query
    select_board(board_name)
    const request = new GetTasksFromBoardRequest()
    request.query = query!
    api.get_tasks_from_board(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                write_messages(res.errors)
                return
            }
            task_infos_map.value[board_name] = res.boards_tasks
            console.log(res.boards_tasks)
        })
}
function close_board(board_name: string) {
    let opened_board_infos_index = -1
    for (let i = 0; i < opened_board_names.value.length; i++) {
        if (opened_board_names.value[i] === board_name) {
            opened_board_infos_index = i
            break
        }
    }
    if (opened_board_infos_index === -1) {
        return
    }
    opened_board_names.value.splice(opened_board_infos_index, 1)
    query_map.value[board_name] = undefined
    if (watching_board_name.value === board_name) {
        watching_board_name.value = null
    }
}
function select_board(board_name: string) {
    watching_board_name.value = board_name
}
async function write_messages(messages: Array<string>) {
    let is_first = true
    for (let i = 0; i < messages.length; i++) {
        const message_ = messages[i]
        await sleep(is_first ? 0 : 5000)
        write_message(message_)
        is_first = false
    }
}
function updated_check_condition(check_state: CheckState) {
    //TODO
}
function updated_sort_type(sort_tyhpe: SortType) {
    //TODO
}
function updated_search_word(word: string) {
    //TODO
}
function updated_boards_by_user() {
    //TODO
}
function clicked_board(board_name: string) {
    let opened = false
    for (let i = 0; i < opened_board_names.value.length; i++) {
        if (opened_board_names.value[i] === board_name) {
            opened = true
            break
        }
    }
    if (!opened) {
        open_board(board_name)
    }
}
function show_add_task_dialog() {
    add_task_dialog_ref.value?.show()
}
function added_task() {
    //TODO
}
function updated_tags_by_user() {
    //TODO
}
function updated_tags(tags: Array<string>) {
    //TODO
}
function copied_task_id(task_info: TaskInfo) {
    //TODO
}
function added_tag() {
    //TODO
}
function added_text() {
    //TODO
}
function updated_task(task_info: TaskInfo) {
    //TODO
}
function deleted_task(task_info: TaskInfo) {
    //TODO
}
function deleted_tag() {
    //TODO
}
function deleted_text() {
    //TODO
}
</script>
<style></style>