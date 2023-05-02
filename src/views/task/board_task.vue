<template>
    <v-card @contextmenu.prevent="show_contextmenu">
        <table>
            <tr>
                <td>
                    <v-checkbox v-model="check" />
                </td>
                <td>
                    {{ title }}
                </td>
                <td v-if="limit">
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
import { Ref, ref, watch } from 'vue';
import task_contextmenu from './task_contextmenu.vue';
import generate_uuid from '@/generate_uuid';

interface Props {
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'copied_task_id', task_info: TaskInfo): void
    (e: 'added_tag'): void
    (e: 'added_text'): void
    (e: 'updated_task', task_info: TaskInfo): void
    (e: 'deleted_task', task_info: TaskInfo): void
}>()

let check: Ref<boolean> = ref(props.task_info.check_state_info.is_checked)
let title: Ref<string> = ref(props.task_info.task_title_info.title)
let limit: Ref<Date | null> = ref(props.task_info.limit_info.limit)
let x_contextmenu: Ref<number> = ref(0)
let y_contextmenu: Ref<number> = ref(0)
const task_context_menu_ref = ref<InstanceType<typeof task_contextmenu> | null>(null);

watch(() => props.task_info, () => {
    check.value = props.task_info.check_state_info.is_checked
    title.value = props.task_info.task_title_info.title
    limit.value = props.task_info.limit_info.limit
    console.log(title.value)
})


watch(check, () => {
    const api = new MiServerAPI()

    const new_task_info = new TaskInfo()
    new_task_info.task = props.task_info.task
    new_task_info.task_title_info = props.task_info.task_title_info
    new_task_info.check_state_info = new CheckStateInfo()
    new_task_info.limit_info = props.task_info.limit_info
    new_task_info.board_info = props.task_info.board_info

    new_task_info.check_state_info.check_state_id = generate_uuid()
    new_task_info.check_state_info.task_id = new_task_info.task.task_id
    new_task_info.check_state_info.updated_time = new Date(Date.now())
    new_task_info.check_state_info.is_checked = check.value

    const update_request = new UpdateTaskRequest()
    update_request.task_info = new_task_info

    api.update_task(update_request)
        .then(res => {
            if (res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emits("updated_task", new_task_info)
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
    emits("added_tag")//TODO 一貫性をもたせるならばAPIでTask取得してEmitしたほうがいい
}
function emit_added_text() {
    emits("added_text")//TODO 一貫性をもたせるならばAPIでText取得してEmitしたほうがいい
}
function emit_updated_task(updated_task_info: TaskInfo) {
    emits("updated_task", updated_task_info)
}
function emit_deleted_task(deleted_task_info: TaskInfo) {
    emits("deleted_task", deleted_task_info)
}
</script>

<style></style>