<template>
    <v-card @click="emit_clicked_board" dropzone="true" @drop="drop" @dragover="dragover">
        <v-card-title :style="title_style">
            <v-container class="pa-0 ma-0">
                <v-row class="pa-0 ma-0">
                    <v-col class="pa-0 ma-0" cols="auto">
                        {{ board_name }}
                    </v-col>
                    <v-spacer />
                    <v-col class="pa-0 ma-0" cols="auto">
                        <v-btn icon="mdi-close" @click="emit_close_board_request" />
                    </v-col>
                </v-row>
            </v-container>
        </v-card-title>
        <board_task v-for="task_info in task_infos" :key="task_info.task.task_id" :task_info="task_info"
            @errors="emit_errors" @copied_task_id="emit_copied_task_id" @added_tag="emit_added_tag"
            @added_text="emit_added_text" @updated_task="emit_updated_task" @deleted_task="emit_deleted_task"
            @clicked_task="emit_clicked_task" />
    </v-card>
</template>

<script setup lang="ts">
import TaskInfo from '@/api/data_struct/TaskInfo';
import board_task from '../task/board_task.vue';
import { Ref, ref, watch, nextTick } from 'vue';
import BoardInfo from '@/api/data_struct/BoardInfo';
import generate_uuid from '@/generate_uuid';
import MiServerAPI from '@/api/MiServerAPI';
import UpdateTaskRequest from '@/api/UpdateTaskRequest';

interface Props {
    task_infos: Array<TaskInfo>
    board_name: string
    selected_board_name: string | null
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'copied_task_id', task_info: TaskInfo): void
    (e: 'added_tag'): void
    (e: 'added_text'): void
    (e: 'updated_task', old_task_info: TaskInfo, new_task_info: TaskInfo): void
    (e: 'deleted_task', task_info: TaskInfo): void
    (e: 'clicked_task', task_info: TaskInfo): void
    (e: 'clicked_board', board_name: string): void
    (e: 'close_board_request', board_name: string): void
}>()

const title_style: Ref<any> = ref(generate_title_style())

watch(() => props.selected_board_name, () => {
    update_style()
})

function update_style() {
    title_style.value = generate_title_style()
}
function generate_title_style(): any {
    return { background: props.selected_board_name == props.board_name ? "whitesmoke" : "white" }
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_copied_task_id(task_info: TaskInfo) {
    emits("copied_task_id", task_info)
}
function emit_added_tag() {
    emits("added_tag")
}
function emit_added_text() {
    emits("added_text")
}
function emit_updated_task(old_task_info: TaskInfo, new_task_info: TaskInfo) {
    emits("updated_task", old_task_info, new_task_info)
}
function emit_deleted_task(deleted_task_info: TaskInfo) {
    emits("deleted_task", deleted_task_info)
}
function emit_clicked_task(clicked_task_info: TaskInfo) {
    emits("clicked_task", clicked_task_info)
}
function emit_close_board_request() {
    emits("close_board_request", props.board_name)
}
function emit_clicked_board() {
    emits("clicked_board", props.board_name)
}
function dragover(e: DragEvent) {
    e!.dataTransfer!.dropEffect = "move"
    e!.preventDefault()
    e!.stopPropagation()
}
function drop(e: DragEvent) {
    let drop_task_info: TaskInfo = new TaskInfo()
    try {
        drop_task_info = JSON.parse(e.dataTransfer!.getData("mi/task_info"))
        if (drop_task_info.task.task_id == "") {
            return
        }
    } catch {
        return
    }
    e!.preventDefault()
    e!.stopPropagation()
    const api = new MiServerAPI()
    const new_task_info = new TaskInfo()
    new_task_info.task = drop_task_info.task
    new_task_info.task_title_info = drop_task_info.task_title_info
    new_task_info.check_state_info = drop_task_info.check_state_info
    new_task_info.limit_info = drop_task_info.limit_info
    new_task_info.board_info = new BoardInfo()
    new_task_info.board_info.board_info_id = generate_uuid()
    new_task_info.board_info.task_id = drop_task_info.task.task_id
    new_task_info.board_info.board_name = props.board_name
    new_task_info.board_info.updated_time = new Date(Date.now())

    const request = new UpdateTaskRequest()
    request.task_info = new_task_info
    api.update_task(request)
        .then((res) => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_updated_task(drop_task_info, new_task_info)
        })
}
</script>

<style></style>