<template>
    <v-dialog v-model="is_show">
        <v-card class="pa-5">
            <v-card-title>
                テキスト追加
            </v-card-title>
            <v-text-area v-model="text" tabindex="101" />
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="close_dialog" tabindex="103">
                            キャンセル
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="submit" tabindex="102">
                            追加
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
    show: boolean
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'added_text'): void
}>()

let text: Ref<string> = ref("")
let is_show: Ref<boolean> = ref(props.show)

watch(() => props.show, () => {
    is_show.value = props.show
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
            if (res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_added_text()
            clear_fields()
            close_dialog()
        })
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