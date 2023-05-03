<template>
    <div class="text" @contextmenu.prevent.stop="show_contextmenu">
        <p class="text_content">{{ text?.text }}</p>
        <text_contextmenu :text="text" :x="x_contextmenu" :y="y_contextmenu" @errors="emit_errors"
            @deleted_text="emit_deleted_text" ref="contextmenu" />
    </div>
</template>

<script setup lang="ts">
import Text from '@/api/data_struct/Text';
import text_contextmenu from './text_contextmenu.vue';
import { Ref, ref } from 'vue';

interface Props {
    text: Text
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'deleted_text'): void
}>()
const contextmenu = ref<InstanceType<typeof text_contextmenu> | null>(null);

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
function emit_deleted_text() {
    emits("deleted_text")
}
</script>

<style>
.text {
  background-color: #eee;
  border: dashed 1px;
  margin: 8px;
  padding: 8px;
}
.text_content {
  white-space: pre-line;
}
</style>