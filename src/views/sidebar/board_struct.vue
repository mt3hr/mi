<template>
    <tr v-if="is_item()">
        <td>
            <table>
                <tr>
                    <td class="tree_item ml-1">{{ struct.key }}</td>
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
                    @updated_check_items_by_user="emit_updated_check_items_by_user"
                    @click_items_by_user="emit_click_items_by_user" />
            </table>
        </td>
    </tr>
</template>

<script setup lang="ts">
</script>

<style></style>