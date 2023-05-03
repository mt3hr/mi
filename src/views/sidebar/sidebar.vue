<template>
    <v-card>
        <check_condition_selectbox @errors="emit_errors" @updated_check_condition="updated_check_condition" />
        <sort_condition_selectbox @errors="emit_errors" @updated_sort_type="updated_sort_type" />
        <word_search_textbox @errors="emit_errors" @updated_search_word="updated_search_word" />
        <board_list :option="option" @errors="emit_errors" @updated_by_user="updated_boards_by_user"
            @clicked_board="clicked_board" />
        <tag_list :option="option" :checked_tags="checked_tags" @errors="emit_errors"
            @updated_by_user="updated_tags_by_user" @updated_checked_tags="updated_tags" />
    </v-card>
</template>

<script setup lang="ts">
import { Ref, ref, watch } from 'vue';
import check_condition_selectbox from './check_condition_selectbox.vue';
import sort_condition_selectbox from './sort_condition_selectbox.vue';
import word_search_textbox from './word_search_textbox.vue';
import board_list from './board_list.vue';
import tag_list from './tag_list.vue';
import CheckState from '@/api/data_struct/CheckState';
import SortType from '@/api/data_struct/SortType';
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
}>()

let option: Ref<ApplicationConfig> = ref(new ApplicationConfig())
let checked_tags: Ref<Array<string>> = ref(new Array<string>())

function updated_check_condition(check_state: CheckState) {
    //TODO
}
function updated_sort_type(sort_type: SortType) {
    //TODO
}
function updated_search_word(word: string) {
    //TODO
}
function updated_boards_by_user() {
    //TODO
}
function clicked_board(board: any) {
    //TODO
}
function updated_tags_by_user() {
    //TODO
}
function updated_tags(tags: Array<string>) {
    //TODO
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_updated_check_condition(check_state: CheckState) {
    //TODO
}
function emit_updated_sort_type(sort_type: SortType) {
    //TODO
}
function emit_updated_search_word(word: string) {
    //TODO
}
function emit_updated_boards_by_user() {
    //TODO
}
function emit_clicked_board(board: any) {
    //TODO
}
function emit_updated_tags_by_user() {
    //TODO
}
function emit_updated_tags(tag: Array<string>) {
    //TODO
}
</script>

<style></style>