<template>
    <v-text-field v-model="search_word" :label="'検索'" @keypress.enter="emit_updated_search_word" />
</template>

<script setup lang="ts">
import { Ref, ref, nextTick, watch } from 'vue';

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_search_word', word: string): void
}>()

let search_word: Ref<string> = ref("")

defineExpose({
    set_search_word_by_application,
    get_search_word
})

function get_search_word(): string {
    return search_word.value
}
function set_search_word_by_application(new_search_word: string): void {
    search_word.value = new_search_word
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_updated_search_word() {
    emits("updated_search_word", search_word.value)
}
</script>

<style></style>