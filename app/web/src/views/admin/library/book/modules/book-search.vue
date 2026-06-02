<script setup lang="ts">
import { toRaw } from 'vue';
import { jsonClone } from '@sa/utils';
import { useNaiveForm } from '@/hooks/common/form';
import { useDictItems } from '@/hooks/business/dict';
import { $t } from '@/locales';

defineOptions({
  name: 'BookSearch'
});

interface Emits {
  (e: 'reset'): void;
  (e: 'search'): void;
}

const emit = defineEmits<Emits>();

const { formRef, validate, restoreValidation } = useNaiveForm();

const model = defineModel<Api.BookManage.BookSearchParams>('model', { required: true });

const { options: serialStatusOptions } = useDictItems('book_serial_status');
const { options: visibilityOptions } = useDictItems('book_visibility');

const defaultModel = jsonClone(toRaw(model.value));

function resetModel() {
  Object.assign(model.value, defaultModel);
}

async function reset() {
  await restoreValidation();
  resetModel();
  emit('reset');
}

async function search() {
  await validate();
  emit('search');
}
</script>

<template>
  <NCard :bordered="false" size="small" class="card-wrapper">
    <NCollapse :default-expanded-names="['book-search']">
      <NCollapseItem :title="$t('common.search')" name="book-search">
        <NForm ref="formRef" :model="model" label-placement="left" :label-width="80">
          <NGrid responsive="screen" item-responsive>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.admin.library.book.bookName')" path="title" class="pr-24px">
              <NInput v-model:value="model.title" :placeholder="$t('page.admin.library.book.form.title')" />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.admin.library.book.author')" path="author" class="pr-24px">
              <NInput v-model:value="model.author" :placeholder="$t('page.admin.library.book.form.author')" />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.admin.library.book.serialStatus')" path="serialStatus" class="pr-24px">
              <NSelect
                v-model:value="model.serialStatus"
                :placeholder="$t('page.admin.library.book.form.serialStatus')"
                :options="serialStatusOptions"
                clearable
              />
            </NFormItemGi>
            <NFormItemGi span="24 s:12 m:6" :label="$t('page.admin.library.book.visibility')" path="visibility" class="pr-24px">
              <NSelect
                v-model:value="model.visibility"
                :placeholder="$t('page.admin.library.book.form.visibility')"
                :options="visibilityOptions"
                clearable
              />
            </NFormItemGi>
            <NFormItemGi span="24 m:12" class="pr-24px">
              <NSpace class="w-full" justify="end">
                <NButton @click="reset">
                  <template #icon>
                    <icon-ic-round-refresh class="text-icon" />
                  </template>
                  {{ $t('common.reset') }}
                </NButton>
                <NButton type="primary" ghost @click="search">
                  <template #icon>
                    <icon-ic-round-search class="text-icon" />
                  </template>
                  {{ $t('common.search') }}
                </NButton>
              </NSpace>
            </NFormItemGi>
          </NGrid>
        </NForm>
      </NCollapseItem>
    </NCollapse>
  </NCard>
</template>

<style scoped></style>
