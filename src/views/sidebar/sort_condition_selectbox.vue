<template>
    <v-select v-model="sort_type" :label="'ソート'" :items="sort_types" item-title="title" item-value="value"
        @update:model-value="emit_updated_sort_type" />
</template>

<script setup lang="ts">
import SortType from '@/api/data_struct/SortType';
import { Ref, ref, nextTick } from 'vue';
class SortTypeSelectModel {
    title: string = ""
    value: SortType = SortType.CreatedTimeDesc
}
const sort_type_select_model = new Array<SortTypeSelectModel>()
const sort_type_created_time_desc = new SortTypeSelectModel()
sort_type_created_time_desc.title = "チェック済みのみ"
sort_type_created_time_desc.value = SortType.CreatedTimeDesc
const sort_type_limit_time_asc = new SortTypeSelectModel()
sort_type_limit_time_asc.title = "全て"
sort_type_limit_time_asc.value = SortType.LimitTimeAsc
sort_type_select_model.push(sort_type_created_time_desc)
sort_type_select_model.push(sort_type_limit_time_asc)

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_sort_type', sort_type: SortType): void
}>()

const sort_types: Ref<Array<SortTypeSelectModel>> = ref(sort_type_select_model)
const sort_type: Ref<SortType> = ref(SortType.CreatedTimeDesc)

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_updated_sort_type() {
    emits("updated_sort_type", sort_type.value)
}
</script>

<style></style>