<template>
    <tr v-if="is_item()">
        <td>
            <table>
                <tr>
                    <td class="tree_item ml-1" @click="click_item_by_user">{{ struct.key }}</td>
                </tr>
            </table>
        </td>
    </tr>
    <tr v-else>
        <td>
            <table>
                <tr>
                    <td>
                        <span v-if="open_group" style="cursor: default" @click="open_group = !open_group">▽</span>
                        <span v-else style="cursor: default" @click="open_group = !open_group">▷</span>
                    </td>
                    <td @click="click_group_by_user">
                        <div class="tree_item">{{ group_name }}</div>
                    </td>
                </tr>
            </table>
            <table class="ml-4">
                <structures v-show="open_group" v-for="(child_struct, index) in struct_list" :open="get_group_open(index)"
                    :key="index" :struct="child_struct" :group_name="get_group_name(index)"
                    @click_items_by_user="emit_click_items_by_user" />
            </table>
        </td>
    </tr>
</template>

<script setup lang="ts">
import structures from './board_struct.vue'
import { Ref, ref, watch } from 'vue';

interface Props {
    struct: any
    group_name: string
    open: boolean
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'click_items_by_user', items: Array<string>): void
}>()

let open_group: Ref<boolean> = ref(props.open)
let struct_list: Ref<any> = ref(new Array<string>())
let indeterminate_group: Ref<boolean> = ref(false)

watch(() => props.open, () => {
    open_group.value = props.open
})
watch(() => props.struct, () => {
    updated_struct()
})

open_group.value = props.open
updated_struct()

// アイテムではなくの場合に使われます。
// 子アイテムを子アイテム配列に変換してthis.struct_listに収めます。
// this.struct_listはv-forで回して子アイテムとして再帰的に読み込まれます。
function updated_struct() {
    struct_list.value = Object.values(props.struct)
}
// this.structがアイテムであればtrueを、そうではなくグループである場合はfalseを返します。
function is_item() {
    return props.struct.key !== undefined
}
function get_group_open(index: number) {
    let group_name = Object.keys(props.struct)[index]
    if (group_name.endsWith(',close') || group_name.endsWith(', close')) {
        return false
    } else if (group_name.endsWith(',open') || group_name.endsWith(', open')) {
        return true
    }
    return true
}
// 子アイテムのグループ名を取得するためにv-for内から使われます。
function get_group_name(index: number) {
    let group_name = Object.keys(props.struct)[index]
    if (group_name.endsWith(',close') || group_name.endsWith(', close')) {
        group_name = group_name.split(',').slice(0, -1).join(',')
    } else if (group_name.endsWith(',open') || group_name.endsWith(', open')) {
        group_name = group_name.split(',').slice(0, -1).join(',')
    }
    return group_name
}
// 子グループ内の複数のアイテムのみをチェックするように変更があったときに、それを上に伝えるために呼び出されます。
function emit_click_items_by_user(items: Array<string>) {
    emits('click_items_by_user', items)
}
// 子グループ内の一つのアイテムのみをチェックするよう変更があったときに、それを上に伝えるために呼び出されます。
function emit_click_item_by_user(item: string) {
    emit_click_items_by_user([item])
}
// このアイテムがクリックされたときに呼び出されます。
// このアイテムのみにチェックが入るように上にemitします。
function click_item_by_user() {
    emit_click_item_by_user(props.struct.key)
}
// このアイテムがクリックされたときに呼び出されます。
// このアイテム内のアイテムのみにチェックが入るように上にemitします。
function click_group_by_user() {
    let items = new Array<string>()
    let f = (struct: any) => { }
    let func = (struct: any) => {
        if (struct.key !== undefined) {
            items.push(struct.key)
        } else {
            Object.keys(struct).forEach(name => {
                f(struct[name])
            })
        }
    }
    f = func
    f(props.struct)
    emit_click_items_by_user(items)
}
</script>

<style>
.tree_item {
    min-width: 200px;
    cursor: default;
}
</style>