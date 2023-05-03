<template>
    <v-card>
        <v-card-title>
            {{ board_name }}
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

interface Props {
    task_infos: Array<TaskInfo>
    board_name: string
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'copied_task_id', task_info: TaskInfo): void
    (e: 'added_tag'): void
    (e: 'added_text'): void
    (e: 'updated_task', task_info: TaskInfo): void
    (e: 'deleted_task', task_info: TaskInfo): void
    (e: 'clicked_task', task_info: TaskInfo): void
}>()

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_copied_task_id(task_info: TaskInfo) {
    emits("copied_task_id", task_info)
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
function emit_clicked_task(clicked_task_info: TaskInfo) {
    emits("clicked_task", clicked_task_info)
}
</script>

<style></style>