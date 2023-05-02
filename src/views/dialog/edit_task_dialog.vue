<template>
</template>

<script setup lang="ts">
import TaskInfo from '@/api/data_struct/TaskInfo';

interface Props {
    task_info: TaskInfo
}
const props = defineProps<Props>()

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_task', task_info: TaskInfo): void
}>()

function emit_updated_task() {
    emits("updated_task", props.task_info) //TODO 古い情報だけどいいの？
}
</script>

<style>
</style>