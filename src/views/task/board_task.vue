<template>
    <v-card :id="task_info.task.task_id" @contextmenu.prevent="show_contextmenu" @click="emit_clicked_task"
        class="board_task_card pa-0 ma-0" :draggable="true" @dragstart="dragstart">
        <table class="task_title_line_table">
            <tr>
                <td align="left" class="task_checkbox_td">
                    <v-checkbox v-model="check" class="pa-0 ma-0" />
                </td>
                <td class="task_title_td pa-0 ma-0" align="left">
                    <p>{{ title }}</p>
                </td>
                <td>
                    <table>
                        <tr>
                            <td v-if="sort_type === SortType.LimitTimeAsc" class="time_td" align="right">
                                <small>
                                    <p v-if="limit">期限: {{ limit.toLocaleString() }}</p>
                                </small>
                            </td>
                        </tr>
                        <tr>
                            <td v-if="sort_type === SortType.StartTimeDesc" class="time_td" align="right">
                                <small>
                                    <p v-if="start">開始: {{ start.toLocaleString() }}</p>
                                </small>
                            </td>
                        </tr>
                        <tr>
                            <td v-if="sort_type === SortType.EndTimeDesc" class="time_td" align="right">
                                <small>
                                    <p v-if="end">終了: {{ end.toLocaleString() }}</p>
                                </small>
                            </td>
                        </tr>
                        <tr>
                            <td v-if="sort_type === SortType.CreatedTimeDesc" class="time_td" align="right">
                                <small>
                                    <p v-if="created_time">作成: {{ created_time.toLocaleString() }}</p>
                                </small>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
        <task_contextmenu :task_info="task_info" :x="x_contextmenu" :y="y_contextmenu" @added_tag="emit_added_tag"
            @added_text="emit_added_text" @copied_task_id="emit_copied_task_id" @updated_task="emit_updated_task"
            @deleted_task="emit_deleted_task" @errors="emit_errors" ref="task_context_menu_ref" />
    </v-card>
</template>

<script setup lang="ts">
import CheckStateInfo from '@/api/data_struct/CheckStateInfo';
import TaskInfo from '@/api/data_struct/TaskInfo'
import MiServerAPI from '@/api/MiServerAPI';
import UpdateTaskRequest from '@/api/UpdateTaskRequest';
import { Ref, ref, watch, nextTick } from 'vue';
import task_contextmenu from './task_contextmenu.vue';
import generate_uuid from '@/generate_uuid';
import SortType from '@/api/data_struct/SortType';

interface Props {
    task_info: TaskInfo
    sort_type: SortType
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
}>()

let check: Ref<boolean> = ref(props.task_info.check_state_info.is_checked)
let title: Ref<string> = ref(props.task_info.task_title_info.title)
let created_time: Ref<Date | null> = ref(props.task_info.task.created_time)
let limit: Ref<Date | null> = ref(props.task_info.limit_info.limit)
let start: Ref<Date | null> = ref(props.task_info.start_info.start)
let end: Ref<Date | null> = ref(props.task_info.end_info.end)
let x_contextmenu: Ref<number> = ref(0)
let y_contextmenu: Ref<number> = ref(0)
const task_context_menu_ref = ref<InstanceType<typeof task_contextmenu> | null>(null);

watch(() => props.task_info, () => {
    check.value = props.task_info.check_state_info.is_checked
    title.value = props.task_info.task_title_info.title
    created_time.value = props.task_info.task.created_time
    limit.value = props.task_info.limit_info.limit
    start.value = props.task_info.start_info.start
    end.value = props.task_info.end_info.end
})


watch(check, () => {
    const api = new MiServerAPI()

    const new_task_info = new TaskInfo()
    new_task_info.task = props.task_info.task
    new_task_info.task_title_info = props.task_info.task_title_info
    new_task_info.check_state_info = new CheckStateInfo()
    new_task_info.limit_info = props.task_info.limit_info
    new_task_info.start_info = props.task_info.start_info
    new_task_info.end_info = props.task_info.end_info
    new_task_info.board_info = props.task_info.board_info

    new_task_info.check_state_info.check_state_id = generate_uuid()
    new_task_info.check_state_info.task_id = new_task_info.task.task_id
    new_task_info.check_state_info.updated_time = new Date(Date.now())
    new_task_info.check_state_info.is_checked = check.value

    const update_request = new UpdateTaskRequest()
    update_request.task_info = new_task_info

    api.update_task(update_request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_updated_task(props.task_info, new_task_info)
        })

})
function show_contextmenu(e: MouseEvent) {
    x_contextmenu.value = e.x
    y_contextmenu.value = e.y
    task_context_menu_ref.value!.show()
}
function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_copied_task_id() {
    emits("copied_task_id", props.task_info)
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
function emit_clicked_task() {
    emits("clicked_task", props.task_info)
}
function dragstart(e: DragEvent) {
    e.dataTransfer!.setData("mi/task_info", JSON.stringify(props.task_info))
}
</script>

<style>
.v-checkbox>.v-input__details {
    height: 0 !important;
    max-height: 0 !important;
    min-height: 0 !important;
}

.task_checkbox_td {
    width: 40px;
    max-width: 40px;
    min-width: 40px;
}

.task_title_td {
    width: 190px;
    max-width: 190px;
    min-width: 190px;
}

.time_td {
    padding-right: 30px;
}

#app {
    overflow-y: hidden;
    max-height: 100vh;
}
</style>