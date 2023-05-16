<template>
    <v-navigation-drawer v-model="show_drawer" app>
        <sidebar :option="option" ref="sidebar_ref" @errors="write_messages"
            @updated_check_condition="updated_check_condition" @updated_sort_type="updated_sort_type"
            @updated_search_word="updated_search_word" @updated_boards_by_user="updated_boards_by_user"
            @clicked_board="clicked_board_at_sidebar" @updated_tags_by_user="updated_tags_by_user"
            @updated_checked_tags="updated_checked_tags" />
    </v-navigation-drawer>

    <v-app-bar class="app_bar" app color="indigo" flat dark height="50px">
        <v-app-bar-nav-icon @click.stop="show_drawer = !show_drawer" />
        <v-toolbar-title>mi</v-toolbar-title>
        <v-spacer />
    </v-app-bar>

    <v-main class="main">
        <div class="boards_view">
            <table class="pa-0 ma-0">
                <tr class="boards_wrap pa-0 ma-0">
                    <td class="board_wrap pa-0 ma-0" cols="auto" v-for="board_name in opened_board_names" :key="board_name">
                        <board class="pa-0 ma-0 board" :board_name="board_name" :selected_board_name="watching_board_name"
                            :task_infos="task_infos_map[board_name]" :loading="loading_map[board_name]"
                            @errors="write_messages" @copied_task_id="copied_task_id" @added_tag="added_tag"
                            @added_text="added_text" @updated_task="updated_task" @deleted_task="deleted_task"
                            @clicked_task="set_watching_task" @close_board_request="close_board"
                            @clicked_board="clicked_board_at_sidebar" />
                    </td>
                </tr>
            </table>
        </div>
        <detail_task class="detail_task pa-0 ma-0" v-if="watching_task_info" :task_info="watching_task_info"
            @copied_task_id="copied_task_id" @added_tag="added_tag" @added_text="added_text" @updated_task="updated_task"
            @deleted_task="deleted_task" @deleted_tag="deleted_tag" @deleted_text="deleted_text" ref="detail_task_ref" />
    </v-main>
    <v-avatar :style="floatingActionButtonStyle()" color="indigo" class="position-fixed">
        <v-btn color="white" icon="mdi-plus" variant="text" @click="show_add_task_dialog" />
    </v-avatar>

    <add_task_dialog :option="option" @errors="write_messages" @added_task="added_task" ref="add_task_dialog_ref" />

    <v-snackbar v-model="show_message_snackbar">
        <v-container class="ma-0 pa-0">
            <v-row class="ma-0 pa-0">
                <v-col cols="11" class="ma-0 pa-0">
                    <p>{{ message }}</p>
                </v-col>
                <v-col cols="1" class="ma-0 pa-0">
                    <v-btn icon="mdi-close" @click="show_message_snackbar = false" width="20px" height="20px"
                        class="ma-0 pa-0" />
                </v-col>
            </v-row>
        </v-container>
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
const detail_task_ref = ref<InstanceType<typeof detail_task> | null>(null);
const query_map: Ref<any> = ref({})
const abort_controller_map: Ref<any> = ref({})
const task_infos_map: Ref<any> = ref({})
const loading_map: Ref<any> = ref({})

const actual_height = window.innerHeight;
const element_height = document!.querySelector('#control-height') ? document!.querySelector('#control-height')!.clientHeight : actual_height
const bar_height = (actual_height - element_height) + "px";

update_option()
    .then(() => open_board(option.value?.default_board_name))

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))
function floatingActionButtonStyle() {
    return {
        'bottom': '10px',
        'right': '10px',
        'height': '50px',
        'width': '50px'
    }
}
async function update_option() {
    const api = new MiServerAPI()
    const request = new GetApplicationConfigRequest()
    await api.get_application_config(request)
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
function set_watching_task(task_info: TaskInfo) {
    watching_task_info.value = task_info
}
function open_board(board_name: string) {
    if (!board_name || board_name === "") {
        return
    }
    opened_board_names.value.push(board_name)
    select_board(board_name)
    update_board(board_name)
}
function update_board(board_name: string) {
    const api = new MiServerAPI()
    const query = sidebar_ref.value?.construct_task_search_query()
    if (!query || board_name == "") {
        return
    }
    loading_map.value[board_name] = true
    if (abort_controller_map.value[board_name]) {
        abort_controller_map.value[board_name].abort()
    }
    const abort_controller = new AbortController()
    abort_controller_map.value[board_name] = abort_controller

    query!.board = board_name
    query_map.value[board_name] = query
    const request = new GetTasksFromBoardRequest()
    request.query = query!
    api.get_tasks_from_board(request, abort_controller)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                write_messages(res.errors)
                loading_map.value[board_name] = false
                return
            }
            task_infos_map.value[board_name] = res.boards_tasks
            loading_map.value[board_name] = false
        })
        .catch((err) => {
            return // DOMException: The user aborted a request.が飛んで邪魔なので握りつぶす
        })
}
function close_board(board_name: string) {
    task_infos_map.value[board_name] = undefined
    let opened_board_names_index = -1
    for (let i = 0; i < opened_board_names.value.length; i++) {
        if (opened_board_names.value[i] === board_name) {
            opened_board_names_index = i
            break
        }
    }
    if (opened_board_names_index === -1) {
        return
    }
    opened_board_names.value.splice(opened_board_names_index, 1)
    query_map.value[board_name] = undefined
    if (watching_board_name.value === board_name) {
        watching_board_name.value = null
    }
}
async function select_board(board_name: string | null) {
    watching_board_name.value = board_name
    if (!board_name || board_name === "") {
        watching_task_info.value = null
        return
    }
    if (!query_map.value[board_name!]) {
        await sidebar_ref.value?.check_all_tags()
        update_board(board_name!)
    }
    sidebar_ref.value?.set_search_word_by_application(query_map.value![board_name!].word)
    sidebar_ref.value?.set_sort_type_by_application(query_map.value![board_name!].sort_type)
    sidebar_ref.value?.set_check_state_by_application(query_map.value![board_name!].check_state)
    sidebar_ref.value?.set_checked_tags_by_application(query_map.value![board_name!].tags)
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
    update_board(watching_board_name.value!)
}
function updated_sort_type(sort_tyhpe: SortType) {
    update_board(watching_board_name.value!)
}
function updated_search_word(word: string) {
    update_board(watching_board_name.value!)
}
function updated_boards_by_user() {
    update_board(watching_board_name.value!)
}
function is_opened_board(board_name: string): boolean {
    let opened = false
    for (let i = 0; i < opened_board_names.value.length; i++) {
        if (opened_board_names.value[i] === board_name) {
            opened = true
            break
        }
    }
    return opened
}
function clicked_board_at_sidebar(board_name: string) {
    if (!is_opened_board(board_name)) {
        open_board(board_name)
    }
    select_board(board_name)
}
function show_add_task_dialog() {
    add_task_dialog_ref.value?.show()
}
function added_task(task_info: TaskInfo) {
    const target_board_name = task_info.board_info.board_name
    if (is_opened_board(target_board_name)) {
        select_board(target_board_name)
        update_board(target_board_name)
    }
    update_board_struct()
    write_message("タスクを追加しました")
}
function updated_tags_by_user() {
    if (watching_board_name.value) {
        update_board(watching_board_name.value!)
    }
}
function updated_checked_tags(tags: Array<string>) {
    return
}
function copied_task_id(task_info: TaskInfo) {
    write_message(`コピーしました`)
}
function added_tag() {
    detail_task_ref.value?.update_tags()
    detail_task_ref.value?.update_texts()
    sidebar_ref.value?.update_tag_struct_promise()
    write_message("タグを追加しました")
}
function added_text() {
    detail_task_ref.value?.update_tags()
    detail_task_ref.value?.update_texts()
    write_message("テキストを追加しました")
}
function updated_task(old_task_info: TaskInfo, new_task_info: TaskInfo) {
    const old_board_name = old_task_info.board_info.board_name
    const new_board_name = new_task_info.board_info.board_name
    if (is_opened_board(old_board_name)) {
        select_board(old_board_name)
        update_board(old_board_name)
    }
    if (is_opened_board(new_board_name)) {
        select_board(new_board_name)
        update_board(new_board_name)
    }
    update_board_struct()
    if (watching_task_info.value?.task.task_id === new_task_info.task?.task_id) {
        watching_task_info.value = new_task_info
    }
    write_message("タスクを更新しました")
}
function deleted_task(task_info: TaskInfo) {
    const target_board_name = task_info.board_info.board_name
    if (is_opened_board(target_board_name)) {
        select_board(target_board_name)
        update_board(target_board_name)
    }
    if (watching_task_info.value?.task.task_id === task_info.task?.task_id) {
        select_board(null)
    }
    write_message("タスクを削除しました")
}
function deleted_tag() {
    detail_task_ref.value?.update_tags()
    detail_task_ref.value?.update_texts()
    nextTick(() => sidebar_ref.value?.update_tag_struct_promise())
    write_message("タグを削除しました")
}
function deleted_text() {
    detail_task_ref.value?.update_tags()
    detail_task_ref.value?.update_texts()
    write_message("テキストを削除しました")
}
function update_board_struct() {
    sidebar_ref.value?.update_board_struct_promise()
}
</script>
<style>
.main {
    padding-top: 50px;
}

.board {
    height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    max-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    min-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    width: 390px;
    min-width: 390px;
    max-width: 390px;
    overflow-y: scroll;
}

.detail_task_card {
    height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    max-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    min-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    overflow-y: scroll;
}

.boards_wrap {
    height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    max-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    min-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
}

.detail_task_row {
    height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    max-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
    min-height: calc((100vh - 50px + v-bind(bar_height)) / 2);
}

.boards_view {
    width: 100vw;
    max-width: 100vw;
    min-width: 100vw;
    overflow-x: scroll;
    -ms-overflow-style: none;
    scrollbar-width: none;
}

.boards_view::-webkit-scrollbar {
    display: none;
}

.html {
    overflow-y: hidden;
}

.board .task_title_line_table {
    width: 370px;
    min-width: 370px;
    max-width: 370px;
}

.app_bar {
    height: 50px;
    max-height: 50px;
    min-height: 50px;
}

#control-height {
    height: 100vh;
    width: 0;
    position: absolute;
}
</style>