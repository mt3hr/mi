<template>
    <v-text-field v-model="search_word" :label="'検索'" @keypress.enter="emit_updated_search_word"/>
</template>

<script setup lang="ts">
import { Ref, ref, nextTick } from 'vue';

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_search_word', word: string): void
}>()

let search_word: Ref<string> = ref("")

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_updated_search_word() {
    emits("updated_search_word", search_word.value)
}
</script>

<style></style>