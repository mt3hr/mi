<template>
    <v-dialog v-model="is_show">
        <v-card class="pa-5">
            <v-card-title>
                新規板名
            </v-card-title>
            <v-text-field v-model="board_name" @keypress.enter="submit" tabindex="101" />
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
                            決定
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { Ref, ref, watch } from 'vue';

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'inputed_board_name', board_name: string): void
}>()

let board_name: Ref<string> = ref("")
let is_show: Ref<boolean> = ref(false)

defineExpose({show})

watch(() => is_show.value, () => {
    is_show.value = is_show.value
})

function close_dialog() {
    is_show.value = false
}
function submit() {
    emit_inputed_board_name()
    clear_fields()
    close_dialog()
}
function show() {
    is_show.value = true
}
function clear_fields() {
    board_name.value = ""
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_inputed_board_name() {
    emits("inputed_board_name", board_name.value)
}
</script>

<style></style>