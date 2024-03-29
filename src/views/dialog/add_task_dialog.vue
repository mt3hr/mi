<template>
    <v-dialog v-model="is_show" :width="500">
        <v-card class="pa-5">
            <v-card-title>
                タスク追加
            </v-card-title>
            <v-text-field v-model="task_title" :label="'タイトル'" :autofocus="true" />
            <v-row>
                <v-col cols="10">
                    <v-select :items="board_names" :label="'板名'" v-model="board_name" />
                </v-col>
                <v-col cols="2">
                    <v-btn icon="mdi-plus" @click="show_input_new_board_name_dialog">+</v-btn>
                    <input_new_board_name_dialog @errors="emit_errors" @inputed_board_name="add_and_set_board_name"
                        ref="input_new_board_name_dialog_ref" />
                </v-col>
            </v-row>
            <v-row>
                <v-col cols="auto">
                    <v-checkbox v-model="has_limit" :label="'期日'" />
                </v-col>
                <v-col cols="auto">
                    <input class="input_date" type="date" v-if="has_limit" v-model="limit_date" />
                    <br>
                    <input class="input_time" type="time" v-if="has_limit" v-model="limit_time" />
                </v-col>
            </v-row>
            <v-row>
                <v-col cols="auto">
                    <v-checkbox v-model="has_start" :label="'開始'" />
                </v-col>
                <v-col cols="auto">
                    <input class="input_date" type="date" v-if="has_start" v-model="start_date" />
                    <br>
                    <input class="input_time" type="time" v-if="has_start" v-model="start_time" />
                </v-col>
            </v-row>
            <v-row>
                <v-col cols="auto">
                    <v-checkbox v-model="has_end" :label="'終了'" />
                </v-col>
                <v-col cols="auto">
                    <input class="input_date" type="date" v-if="has_end" v-model="end_date" />
                    <br>
                    <input class="input_time" type="time" v-if="has_end" v-model="end_time" />
                </v-col>
            </v-row>
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="close_dialog">
                            キャンセル
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="submit">
                            追加
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import AddTaskRequest from '@/api/AddTaskRequest';
import MiServerAPI from '@/api/MiServerAPI';
import TaskInfo from '@/api/data_struct/TaskInfo';
import { Ref, ref, watch } from 'vue';
import '@vuepic/vue-datepicker/dist/main.css';
import GetBoardNamesRequest from '@/api/GetBoardNamesRequest';
import input_new_board_name_dialog from './input_new_board_name_dialog.vue';
import generate_uuid from '@/generate_uuid';
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';

interface Props {
    option: ApplicationConfig
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'added_task', task_info: TaskInfo): void
}>()

const now = new Date(Date.now())
const board_names: Ref<Array<string>> = ref(new Array<string>())
const task_title: Ref<string> = ref("")
const board_name: Ref<string> = ref(props.option!.default_board_name)
const is_show: Ref<boolean> = ref(false)
const has_limit: Ref<boolean> = ref(false)
const limit_date: Ref<string> = ref(`${now.getFullYear().toString().padStart(4, '0')}-${(now.getMonth() + 1).toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`)
const limit_time: Ref<string> = ref(`00:00:00`)
const has_start: Ref<boolean> = ref(false)
const start_date: Ref<string> = ref(`${now.getFullYear().toString().padStart(4, '0')}-${(now.getMonth() + 1).toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`)
const start_time: Ref<string> = ref(`00:00:00`)
const has_end: Ref<boolean> = ref(false)
const end_date: Ref<string> = ref(`${now.getFullYear().toString().padStart(4, '0')}-${(now.getMonth() + 1).toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`)
const end_time: Ref<string> = ref(`00:00:00`)

const input_new_board_name_dialog_ref = ref<InstanceType<typeof input_new_board_name_dialog> | null>(null);

defineExpose({ show })

update_board_names()

watch(() => props.option, () => {
    board_name.value = props.option!.default_board_name
})

async function update_board_names() {
    const api = new MiServerAPI()
    const request = new GetBoardNamesRequest()
    await api.get_board_names(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            board_names.value = res.board_names
        })
}
function close_dialog() {
    is_show.value = false
}
function submit() {
    if (task_title.value == "") {
        return
    }
    if (board_name.value == "") {
        return
    }
    const api = new MiServerAPI()
    const request = new AddTaskRequest()
    request.task_info = construct_task_info()
    api.add_task(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_added_task(request.task_info)
            clear_fields()
            close_dialog()
        })
}
function show() {
    board_name.value = props.option.default_board_name
    update_board_names()
        .then(() => is_show.value = true)
}
function clear_fields() {
    const now = new Date(Date.now())
    task_title.value = ""
    has_limit.value = false
    has_start.value = false
    has_end.value = false
    limit_date.value = `${now.getFullYear().toString().padStart(4, '0')}-${(now.getMonth() + 1).toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`
    limit_time.value = `00:00:00`
    start_date.value = `${now.getFullYear().toString().padStart(4, '0')}-${(now.getMonth() + 1).toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`
    start_time.value = `00:00:00`
    end_date.value = `${now.getFullYear().toString().padStart(4, '0')}-${(now.getMonth() + 1).toString().padStart(2, '0')}-${now.getDate().toString().padStart(2, '0')}`
    end_time.value = `00:00:00`
}
function show_input_new_board_name_dialog() {
    input_new_board_name_dialog_ref.value?.show()
}
function add_and_set_board_name(inputed_board_name: string) {
    board_names.value.push(inputed_board_name)
    board_name.value = inputed_board_name
}
function construct_task_info() {
    const task_id = generate_uuid()
    const now = new Date(Date.now())
    const new_task_info = new TaskInfo()
    new_task_info.task.created_time = now
    new_task_info.task.task_id = task_id
    new_task_info.task_title_info.task_title_id = generate_uuid()
    new_task_info.task_title_info.task_id = task_id
    new_task_info.task_title_info.updated_time = now
    new_task_info.task_title_info.title = task_title.value
    new_task_info.check_state_info.check_state_id = generate_uuid()
    new_task_info.check_state_info.task_id = task_id
    new_task_info.check_state_info.updated_time = now
    new_task_info.check_state_info.is_checked = false
    new_task_info.limit_info.limit_id = generate_uuid()
    new_task_info.limit_info.task_id = task_id
    new_task_info.limit_info.updated_time = now
    new_task_info.start_info.start_id = generate_uuid()
    new_task_info.start_info.task_id = task_id
    new_task_info.start_info.updated_time = now
    new_task_info.end_info.end_id = generate_uuid()
    new_task_info.end_info.task_id = task_id
    new_task_info.end_info.updated_time = now
    if (has_limit.value) {
        new_task_info.limit_info.limit = new Date(Date.parse(`${limit_date.value} ${limit_time.value}:00`))
    } else {
        new_task_info.limit_info.limit = null
    }
    if (has_start.value) {
        new_task_info.start_info.start = new Date(Date.parse(`${start_date.value} ${start_time.value}:00`))
    } else {
        new_task_info.start_info.start = null
    }
    if (has_end.value) {
        new_task_info.end_info.end = new Date(Date.parse(`${end_date.value} ${end_time.value}:00`))
    } else {
        new_task_info.end_info.end = null
    }
    new_task_info.board_info.board_info_id = generate_uuid()
    new_task_info.board_info.task_id = task_id
    new_task_info.board_info.updated_time = now
    new_task_info.board_info.board_name = board_name.value
    return new_task_info
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_added_task(added_task_info: TaskInfo) {
    emits("added_task", added_task_info)
}
</script>

<style>
.input_date {
    width: 120px;
}

.input_time {
    width: 120px;
}
</style>