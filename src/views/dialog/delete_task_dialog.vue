<template>
    <v-dialog v-model="is_show">
        <v-card class="pa-5">
            <v-card-title>
                タスク削除
            </v-card-title>
            <v-card-text>{{ task_info.task_title_info.title }}</v-card-text>
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="close_dialog" tabindex="103">
                            キャンセル
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="delete_task" tabindex="102">
                            削除
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import DeleteTaskRequest from '@/api/DeleteTaskRequest';
import MiServerAPI from '@/api/MiServerAPI';
import TaskInfo from '@/api/data_struct/TaskInfo';
import { Ref, ref, watch } from 'vue';

interface Props {
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'deleted_task', task_info: TaskInfo): void
}>()

let is_show: Ref<boolean> = ref(false)

defineExpose({ show })

watch(() => is_show.value, () => {
    is_show.value = is_show.value
})

function close_dialog() {
    is_show.value = false
}
function delete_task() {
    const api = new MiServerAPI()
    const request = new DeleteTaskRequest()
    request.task_id = props.task_info.task.task_id
    api.delete_task(request)
        .then(res => {
            if (res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_deleted_task(props.task_info)
            close_dialog()
        })
}
function show() {
    is_show.value = true
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_deleted_task(deleted_task_info: TaskInfo) {
    emits("deleted_task", deleted_task_info)
}
</script>

<style></style>