<template>
    <span class="tag" @contextmenu.prevent.stop="show_contextmenu">
        <p>{{ tag?.tag }}</p>
        <tag_contextmenu :tag="tag" :x="x_contextmenu" :y="y_contextmenu" @errors="emit_errors"
            @deleted_tag="emit_deleted_tag" ref="contextmenu" />
    </span>
</template>

<script setup lang="ts">
import Tag from '@/api/data_struct/Tag';
import tag_contextmenu from './tag_contextmenu.vue';
import { Ref, ref } from 'vue';

interface Props {
    tag: Tag
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'deleted_tag'): void
}>()
const contextmenu = ref<InstanceType<typeof tag_contextmenu> | null>(null);

let x_contextmenu: Ref<number> = ref(0)
let y_contextmenu: Ref<number> = ref(0)

function show_contextmenu(e: MouseEvent) {
    x_contextmenu.value = e.x
    y_contextmenu.value = e.y
    contextmenu.value!.show()
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_deleted_tag() {
    emits("deleted_tag")
}
</script>

<style>
.tag {
  border: solid white 2px;
  border-left: 0px;
  color: blue;
  cursor : pointer;
  padding: 0 6px 0 2px;
  font-size: small;
  border-radius: 0 1em 1em 0;
  background: lightgray;
  display: inline-flex;
}
.tag::before {
  content: "・";
  color: white;
}
</style>