<script setup lang="ts">
import { ref } from "vue";
import { NButton, NForm, NFormItem, NInput, NModal } from "naive-ui";
import { $t } from "@/locales";

defineOptions({ name: "BookshelfOperateModal" });

const visible = defineModel<boolean>("visible", { default: false });
const emit = defineEmits<{
  submitted: [groupName: string];
}>();

const groupName = ref($t("page.bookshelf.defaultGroup"));

function handleSubmit() {
  emit("submitted", groupName.value);
  groupName.value = $t("page.bookshelf.defaultGroup");
  visible.value = false;
}
</script>

<template>
  <NModal
    v-model:show="visible"
    :title="$t('page.bookshelf.addToBookshelf')"
    preset="card"
    class="w-360px"
  >
    <NForm>
      <NFormItem :label="$t('page.bookshelf.groupName')">
        <NInput
          v-model:value="groupName"
          :placeholder="$t('page.bookshelf.groupPlaceholder')"
        />
      </NFormItem>
      <div class="flex justify-end gap-8px mt-16px">
        <NButton @click="visible = false">{{ $t("common.cancel") }}</NButton>
        <NButton type="primary" @click="handleSubmit">{{ $t("common.confirm") }}</NButton>
      </div>
    </NForm>
  </NModal>
</template>
