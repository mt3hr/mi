<template>
    <v-dialog v-model="is_show">
        <v-card class="pa-5">
            <v-card-title>
                タスク追加
            </v-card-title>
            <v-text-field v-model="task_title" tabindex="101" />
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="close_dialog" tabindex="113">
                            キャンセル
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="submit" tabindex="112">
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

interface Props {
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'added_task', task_info: TaskInfo): void
}>()

let task_title: Ref<string> = ref("")
let is_show: Ref<boolean> = ref(false)

defineExpose({ show })

watch(() => is_show.value, () => {
    is_show.value = is_show.value
})

function close_dialog() {
    is_show.value = false
}
function submit() {
    if (task_title.value == "") {
        return
    }
    const api = new MiServerAPI()
    const request = new AddTaskRequest()
    //TODO task情報からリクエスト組み立て
    api.add_task(request)
        .then(res => {
            if (res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_added_task(request.task_info)
            clear_fields()
            close_dialog()
        })
}
function show() {
    is_show.value = true
}
function clear_fields() {
    //TODO すべての入力情報を消す
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_added_task(added_task_info: TaskInfo) {
    emits("added_task", added_task_info)
}
</script>

<style></style>