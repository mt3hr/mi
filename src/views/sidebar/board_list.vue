<template>
    <h2>板</h2>
    <table class="boardlist">
        <board_struct ref="board_struct_ref" :group_name="''" :struct="board_structure" :open="true"
            @click_items_by_user="(clicked_items) => { emit_clicked_board(clicked_items[0]) }" />
    </table>
</template>
<script setup lang="ts">
import { Ref, ref, watch, nextTick } from 'vue';
import MiServerAPI from '@/api/MiServerAPI';
import board_struct from './board_struct.vue';
import GetBoardNamesRequest from '@/api/GetBoardNamesRequest';
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';

interface Props {
    option: ApplicationConfig
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_by_user'): void
    (e: 'clicked_board', board: string): void
}>()

let boards: Ref<any> = ref({})
let board_structure: Ref<any> = ref({})
const board_struct_ref = ref<InstanceType<typeof board_struct> | null>(null);
const selected_board: Ref<string> = ref("");

defineExpose({
    set_selected_board_by_application,
    get_selected_board
})

nextTick(() => {
    update_boards_promise()
        .then(() => { return update_board_struct_promise() })
        .then(() => emits('updated_by_user'))
})
// board_structをkv_boardlist_boardsの取り扱える形に変換し、更新します。
function update_board_struct_promise() {
    return new Promise(resolve => { return resolve(null) })
        .then(() => {
            // board structを変換してから収めます。
            /*
            {
              "Inbox": "board",
              "プライベート": {
                "開発": "board",
                "nlog未": "board",
              }
            }
            から
            {
              Inbox: {board: "board"},
              プライベート: {
                  "開発": {board: "開発"},
                  "nlog未": {board: "nlog未"},
              },
            }
            に。
            */
            let structed_boards: any = []
            let f = (board_or_parent: any, board_name: string) => { }
            let func = (board_or_parent: any, boardname: string) => {
                if (board_or_parent === 'board') {
                    for (let i = 0; i < boards.value.length; i++) {
                        if (boards.value[i].board == boardname) {
                            break
                        }
                    }
                    structed_boards.push(boardname)
                    return {
                        key: boardname,
                        indeterminate: false,
                    }
                } else {
                    let board_struct: any = {}
                    Object.keys(board_or_parent).forEach(parent => {
                        board_struct[parent] = f(board_or_parent[parent], parent)
                    })
                    return board_struct
                }
            }
            f = func
            let board_struct_
            if (props.option.board_struct) {
                board_struct_ = f(props.option.board_struct, "")
            } else {
                board_struct_ = f({}, "")
            }

            boards.value.forEach((board: any) => {
                let exist = false
                for (let i = 0; i < structed_boards.length; i++) {
                    if (board.board == structed_boards[i]) {
                        exist = true
                        break
                    }
                }
                if (!exist) {
                    board_struct_[board.board] = {
                        key: board.board,
                        indeterminate: false,
                    }
                }

            })
            structed_boards.forEach((board: any) => {
                let exist = false
                for (let i = 0; i < boards.value.length; i++) {
                    if (board == boards.value[i].board) {
                        exist = true
                        break
                    }
                }
                if (!exist) {
                    boards.value.push({
                        board: board,
                    })
                }
            })
            board_structure.value = board_struct_
        })
}
// タグを最新の状態に更新します。
// タグの選択はすべてfalseに初期化されます。
function update_boards_promise() {
    let api = new MiServerAPI()
    return api.get_board_names(new GetBoardNamesRequest())
        .then((res) => {
            let boardsTemp: any = []
            res.board_names.forEach((board: any) => {
                let t = {
                    board: board,
                }
                boardsTemp.push(t)
            })
            boards.value = boardsTemp
        })
        .then(() => { return update_board_struct_promise() })
}
function get_selected_board(): string {
    return selected_board.value
}
function set_selected_board_by_application(new_selected_board: string): void {
    selected_board.value = new_selected_board
}

function emit_clicked_board(board: string) {
    emits("clicked_board", board)
}
</script>

<style></style>