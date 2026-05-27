import { transformRecordToOption } from '@/utils/common';

export const enableStatusRecord: Record<Api.Common.EnableStatus, App.I18n.I18nKey> = {
  '1': 'common.enable',
  '2': 'common.disable'
};

export const enableStatusOptions = transformRecordToOption(enableStatusRecord);

export const userGenderRecord: Record<Api.SystemManage.UserGender, App.I18n.I18nKey> = {
  '1': 'page.manage.user.gender.male',
  '2': 'page.manage.user.gender.female'
};

export const userGenderOptions = transformRecordToOption(userGenderRecord);

export const menuTypeRecord: Record<Api.SystemManage.MenuType, App.I18n.I18nKey> = {
  '1': 'page.manage.menu.type.directory',
  '2': 'page.manage.menu.type.menu'
};

export const menuTypeOptions = transformRecordToOption(menuTypeRecord);

export const menuIconTypeRecord: Record<Api.SystemManage.IconType, App.I18n.I18nKey> = {
  '1': 'page.manage.menu.iconType.iconify',
  '2': 'page.manage.menu.iconType.local'
};

export const menuIconTypeOptions = transformRecordToOption(menuIconTypeRecord);

export const dataScopeRecord: Record<Api.SystemManage.DataScope, App.I18n.I18nKey> = {
  '1': 'page.manage.role.dataScope.all',
  '2': 'page.manage.role.dataScope.custom',
  '3': 'page.manage.role.dataScope.dept',
  '4': 'page.manage.role.dataScope.deptAndSub',
  '5': 'page.manage.role.dataScope.self'
};

export const dataScopeOptions = transformRecordToOption(dataScopeRecord);
