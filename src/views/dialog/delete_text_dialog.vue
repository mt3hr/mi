<template>
    <v-dialog v-model="is_show">
        <v-card class="pa-5">
            <v-card-title>
                テキスト削除
            </v-card-title>
            <v-card-text>{{ text.text }}</v-card-text>
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="close_dialog" tabindex="103">
                            キャンセル
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="delete_text" tabindex="102">
                            削除
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import DeleteTextRequest from '@/api/DeleteTextRequest';
import MiServerAPI from '@/api/MiServerAPI';
import Text from '@/api/data_struct/Text';
import TaskInfo from '@/api/data_struct/TaskInfo';
import { Ref, ref, watch } from 'vue';

interface Props {
    text: Text 
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'deleted_text'): void
}>()

let is_show: Ref<boolean> = ref(false)

defineExpose({show})

watch(() => is_show.value, () => {
    is_show.value = is_show.value
})

function close_dialog() {
    is_show.value = false
}
function delete_text() {
    const api = new MiServerAPI()
    const request = new DeleteTextRequest()
    request.text_id = props.text.id
    api.delete_text(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_deleted_text()
            close_dialog()
        })
}
function show() {
    is_show.value = true
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_deleted_text() {
    emits("deleted_text")
}
</script>

<style>
</style>