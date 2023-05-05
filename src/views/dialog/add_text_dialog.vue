<template>
    <v-dialog v-model="is_show" :width="500">
        <v-card class="pa-5">
            <v-card-title>
                テキスト追加
            </v-card-title>
            <v-textarea v-model="text" :autofocus="true" />
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="submit">
                            追加
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="close_dialog">
                            キャンセル
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import AddTextRequest from '@/api/AddTextRequest';
import MiServerAPI from '@/api/MiServerAPI';
import TaskInfo from '@/api/data_struct/TaskInfo';
import { Ref, ref, watch } from 'vue';

interface Props {
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'added_text'): void
}>()

let text: Ref<string> = ref("")
let is_show: Ref<boolean> = ref(false)

defineExpose({ show })

watch(() => is_show.value, () => {
    is_show.value = is_show.value
})

function close_dialog() {
    is_show.value = false
}
function submit() {
    if (text.value == "") {
        return
    }
    const api = new MiServerAPI()
    const request = new AddTextRequest()
    request.text = text.value
    request.task_id = props.task_info.task.task_id
    api.add_text(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_added_text()
            clear_fields()
            close_dialog()
        })
}
function show() {
    is_show.value = true
}
function clear_fields() {
    text.value = ""
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_added_text() {
    emits("added_text")
}
</script>

<style></style>