<script setup lang="ts">
import { ref } from "vue";
import { NButton, NForm, NFormItem, NInput, NModal } from "naive-ui";
import { $t } from "@/locales";

export interface BookshelfOperateModalProps {
  visible: boolean;
}

const props = defineProps<BookshelfOperateModalProps>();
const emit = defineEmits<{
  "update:visible": [visible: boolean];
  submitted: [groupName: string];
}>();

const groupName = ref("默认");

function handleClose() {
  emit("update:visible", false);
}

function handleSubmit() {
  emit("submitted", groupName.value);
  groupName.value = "默认";
  handleClose();
}
</script>

<template>
  <NModal
    :show="props.visible"
    :title="$t('page.bookshelf.addToBookshelf')"
    preset="card"
    class="w-360px"
    @close="handleClose"
    @update:show="(val: boolean) => !val && handleClose()"
  >
    <NForm>
      <NFormItem :label="$t('page.bookshelf.groupName')">
        <NInput
          v-model:value="groupName"
          :placeholder="$t('page.bookshelf.groupPlaceholder')"
        />
      </NFormItem>
      <div class="flex justify-end gap-8px mt-16px">
        <NButton @click="handleClose">{{ $t("common.cancel") }}</NButton>
        <NButton type="primary" @click="handleSubmit">{{ $t("common.confirm") }}</NButton>
      </div>
    </NForm>
  </NModal>
</template>
