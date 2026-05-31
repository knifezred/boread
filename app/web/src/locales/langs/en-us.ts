const local: App.I18n.Schema = {
  system: {
    title: 'SoybeanAdmin',
    updateTitle: 'System Version Update Notification',
    updateContent: 'A new version of the system has been detected. Do you want to refresh the page immediately?',
    updateConfirm: 'Refresh immediately',
    updateCancel: 'Later'
  },
  common: {
    action: 'Action',
    add: 'Add',
    addSuccess: 'Add Success',
    backToHome: 'Back to home',
    batchDelete: 'Batch Delete',
    cancel: 'Cancel',
    close: 'Close',
    check: 'Check',
    selectAll: 'Select All',
    expandColumn: 'Expand Column',
    columnSetting: 'Column Setting',
    config: 'Config',
    confirm: 'Confirm',
    delete: 'Delete',
    deleteSuccess: 'Delete Success',
    confirmDelete: 'Are you sure you want to delete?',
    edit: 'Edit',
    warning: 'Warning',
    error: 'Error',
    index: 'Index',
    keywordSearch: 'Please enter keyword',
    logout: 'Logout',
    logoutConfirm: 'Are you sure you want to log out?',
    lookForward: 'Coming soon',
    modify: 'Modify',
    modifySuccess: 'Modify Success',
    noData: 'No Data',
    operate: 'Operate',
    enable: 'Enable',
    disable: 'Disable',
    operateSuccess: 'Operation successful',
    operateFail: "Operation failed",
    createTime: 'Create Time',
    updateTime: 'Update Time',
    pleaseCheckValue: 'Please check whether the value is valid',
    refresh: 'Refresh',
    reset: 'Reset',
    search: 'Search',
    switch: 'Switch',
    tip: 'Tip',
    trigger: 'Trigger',
    update: 'Update',
    updateSuccess: 'Update Success',
    userCenter: 'User Center',
    yesOrNo: {
      yes: 'Yes',
      no: 'No'
    },
  },
  request: {
    logout: 'Logout user after request failed',
    logoutMsg: 'User status is invalid, please log in again',
    logoutWithModal: 'Pop up modal after request failed and then log out user',
    logoutWithModalMsg: 'User status is invalid, please log in again',
    refreshToken: 'The requested token has expired, refresh the token',
    tokenExpired: 'The requested token has expired'
  },
  theme: {
    themeDrawerTitle: 'Theme Configuration',
    tabs: {
      appearance: 'Appearance',
      layout: 'Layout',
      general: 'General',
      preset: 'Preset'
    },
    appearance: {
      themeSchema: {
        title: 'Theme Schema',
        light: 'Light',
        dark: 'Dark',
        auto: 'Follow System'
      },
      grayscale: 'Grayscale',
      colourWeakness: 'Colour Weakness',
      themeColor: {
        title: 'Theme Color',
        primary: 'Primary',
        info: 'Info',
        success: 'Success',
        warning: 'Warning',
        error: 'Error',
        followPrimary: 'Follow Primary'
      },
      themeRadius: {
        title: 'Theme Radius'
      },
      recommendColor: 'Apply Recommended Color Algorithm',
      recommendColorDesc: 'The recommended color algorithm refers to',
      preset: {
        title: 'Theme Presets',
        apply: 'Apply',
        applySuccess: 'Preset applied successfully',
        default: {
          name: 'Default Preset',
          desc: 'Default theme preset with balanced settings'
        },
        dark: {
          name: 'Dark Preset',
          desc: 'Dark theme preset for night time usage'
        },
        compact: {
          name: 'Compact Preset',
          desc: 'Compact layout preset for small screens'
        },
        azir: {
          name: "Azir's Preset",
          desc: 'It is a cold and elegant preset that Azir likes'
        }
      }
    },
    layout: {
      layoutMode: {
        title: 'Layout Mode',
        vertical: 'Vertical Mode',
        horizontal: 'Horizontal Mode',
        'vertical-mix': 'Vertical Mix Mode',
        'vertical-hybrid-header-first': 'Left Hybrid Header-First',
        'top-hybrid-sidebar-first': 'Top-Hybrid Sidebar-First',
        'top-hybrid-header-first': 'Top-Hybrid Header-First',
        vertical_detail: 'Vertical menu layout, with the menu on the left and content on the right.',
        'vertical-mix_detail':
          'Vertical mix-menu layout, with the primary menu on the dark left side and the secondary menu on the lighter left side.',
        'vertical-hybrid-header-first_detail':
          'Left hybrid layout, with the primary menu at the top, the secondary menu on the dark left side, and the tertiary menu on the lighter left side.',
        horizontal_detail: 'Horizontal menu layout, with the menu at the top and content below.',
        'top-hybrid-sidebar-first_detail':
          'Top hybrid layout, with the primary menu on the left and the secondary menu at the top.',
        'top-hybrid-header-first_detail':
          'Top hybrid layout, with the primary menu at the top and the secondary menu on the left.'
      },
      tab: {
        title: 'Tab Settings',
        visible: 'Tab Visible',
        cache: 'Tag Bar Info Cache',
        cacheTip: 'Keep the tab bar information after leaving the page',
        height: 'Tab Height',
        mode: {
          title: 'Tab Mode',
          slider: 'Slider',
          chrome: 'Chrome',
          button: 'Button'
        },
        closeByMiddleClick: 'Close Tab by Middle Click',
        closeByMiddleClickTip: 'Enable closing tabs by clicking with the middle mouse button'
      },
      header: {
        title: 'Header Settings',
        height: 'Header Height',
        breadcrumb: {
          visible: 'Breadcrumb Visible',
          showIcon: 'Breadcrumb Icon Visible'
        }
      },
      sider: {
        title: 'Sider Settings',
        inverted: 'Dark Sider',
        width: 'Sider Width',
        collapsedWidth: 'Sider Collapsed Width',
        mixWidth: 'Mix Sider Width',
        mixCollapsedWidth: 'Mix Sider Collapse Width',
        mixChildMenuWidth: 'Mix Child Menu Width',
        autoSelectFirstMenu: 'Auto Select First Submenu',
        autoSelectFirstMenuTip:
          'When a first-level menu is clicked, the first submenu is automatically selected and navigated to the deepest level'
      },
      footer: {
        title: 'Footer Settings',
        visible: 'Footer Visible',
        fixed: 'Fixed Footer',
        height: 'Footer Height',
        right: 'Right Footer'
      },
      content: {
        title: 'Content Area Settings',
        scrollMode: {
          title: 'Scroll Mode',
          tip: 'The theme scroll only scrolls the main part, the outer scroll can carry the header and footer together',
          wrapper: 'Wrapper',
          content: 'Content'
        },
        page: {
          animate: 'Page Animate',
          mode: {
            title: 'Page Animate Mode',
            fade: 'Fade',
            'fade-slide': 'Slide',
            'fade-bottom': 'Fade Zoom',
            'fade-scale': 'Fade Scale',
            'zoom-fade': 'Zoom Fade',
            'zoom-out': 'Zoom Out',
            none: 'None'
          }
        },
        fixedHeaderAndTab: 'Fixed Header And Tab'
      }
    },
    general: {
      title: 'General Settings',
      watermark: {
        title: 'Watermark Settings',
        visible: 'Watermark Full Screen Visible',
        text: 'Custom Watermark Text',
        enableUserName: 'Enable User Name Watermark',
        enableTime: 'Show Current Time',
        timeFormat: 'Time Format'
      },
      multilingual: {
        title: 'Multilingual Settings',
        visible: 'Display multilingual button'
      },
      globalSearch: {
        title: 'Global Search Settings',
        visible: 'Display GlobalSearch button'
      }
    },
    configOperation: {
      copyConfig: 'Copy Config',
      copySuccessMsg: 'Copy Success, Please replace the variable "themeSettings" in "src/theme/settings.ts"',
      resetConfig: 'Reset Config',
      resetSuccessMsg: 'Reset Success'
    }
  },
  route: {
    login: 'Login',
    403: 'No Permission',
    404: 'Page Not Found',
    500: 'Server Error',
    'iframe-page': 'Iframe',
    admin_dashboard: 'Dashboard',
    admin: 'System',
    admin_system: 'System',
    admin_system_dict: 'Dict Management',
    admin_system_log: 'Log Management',
    admin_system_menu: 'Menu Management',
    admin_system_role: 'Role Management',
    admin_system_dept: 'Department Management',
    admin_system_user: 'User Management',
    "admin_system_user-detail": 'User Detail',
    admin_library: 'Book Management',
    "admin_library_book-category": 'Book Categories',
    "admin_library_book-tag": 'Book Tags',
    "admin_library_book": 'Book Management',
    "admin_library_book-chapter-rule": 'Chapter Rules',
    "admin_library_book-filter-rule": 'Filter Rules',
    home: "Home",
    "book-detail": "Book Detail",
    "book-reader": "Book Reader",
    system: "System Management",
    library: "Library Management"
  },
  page: {
    login: {
      common: {
        loginOrRegister: 'Login / Register',
        userNamePlaceholder: 'Please enter user name',
        phonePlaceholder: 'Please enter phone number',
        codePlaceholder: 'Please enter verification code',
        passwordPlaceholder: 'Please enter password',
        confirmPasswordPlaceholder: 'Please enter password again',
        codeLogin: 'Verification code login',
        confirm: 'Confirm',
        back: 'Back',
        validateSuccess: 'Verification passed',
        loginSuccess: 'Login successfully',
        welcomeBack: 'Welcome back, {userName} !'
      },
      pwdLogin: {
        title: 'Password Login',
        rememberMe: 'Remember me',
        forgetPassword: 'Forget password?',
        register: 'Register',
        otherAccountLogin: 'Other Account Login',
        otherLoginMode: 'Other Login Mode',
        superAdmin: 'Super Admin',
        admin: 'Admin',
        user: 'User'
      },
      codeLogin: {
        title: 'Verification Code Login',
        getCode: 'Get verification code',
        reGetCode: 'Reacquire after {time}s',
        sendCodeSuccess: 'Verification code sent successfully',
        imageCodePlaceholder: 'Please enter image verification code'
      },
      register: {
        title: 'Register',
        agreement: 'I have read and agree to',
        protocol: '《User Agreement》',
        policy: '《Privacy Policy》'
      },
      resetPwd: {
        title: 'Reset Password'
      },
      bindWeChat: {
        title: 'Bind WeChat'
      }
    },
    home: {
      branchDesc:
        'For the convenience of everyone in developing and updating the merge, we have streamlined the code of the main branch, only retaining the homepage menu, and the rest of the content has been moved to the example branch for maintenance. The preview address displays the content of the example branch.',
      greeting: 'Good morning, {userName}, today is another day full of vitality!',
      weatherDesc: 'Today is cloudy to clear, 20℃ - 25℃!',
      projectCount: 'Project Count',
      todo: 'Todo',
      message: 'Message',
      downloadCount: 'Download Count',
      registerCount: 'Register Count',
      schedule: 'Work and rest Schedule',
      study: 'Study',
      work: 'Work',
      rest: 'Rest',
      entertainment: 'Entertainment',
      visitCount: 'Visit Count',
      turnover: 'Turnover',
      dealCount: 'Deal Count',
      projectNews: {
        title: 'Project News',
        moreNews: 'More News',
        desc1: 'Soybean created the open source project soybean-admin on May 28, 2021!',
        desc2: 'Yanbowe submitted a bug to soybean-admin, the multi-tab bar will not adapt.',
        desc3: 'Soybean is ready to do sufficient preparation for the release of soybean-admin!',
        desc4: 'Soybean is busy writing project documentation for soybean-admin!',
        desc5: 'Soybean just wrote some of the workbench pages casually, and it was enough to see!'
      },
      creativity: 'Creativity'
    },
    book: {
      filter: {
        category: 'Category',
        serialStatus: 'Serial Status',
        wordCount: 'Word Count',
        tags: 'Tags',
        updateTime: 'Update Time'
      },
      home: {
        all: 'All',
        searchPlaceholder: 'Search title, author',
        importBooks: 'Import Books',
        relatedWorks: 'Total {total} related works',
        noContent: 'No Content',
        uncategorized: 'Uncategorized',
        noIntro: 'No introduction',
        latestChapter: 'Latest Chapter',
        oneWeek: '1 Week',
        oneMonth: '1 Month',
        threeMonths: '3 Months',
        oneYear: '1 Year',
        sortPopular: 'Popular',
        sortCollect: 'Collection',
        sortWord: 'Words',
        sortVote: 'Votes',
        sortMonthly: 'Monthly'
      },
      detail: {
        breadcrumbHome: 'Home',
        bookInfo: 'Book Info',
        catalog: 'Catalog',
        author: 'Author',
        updateTime: 'Update Time',
        latestChapter: 'Latest Chapter',
        words: 'words',
        chapters: 'chapters',
        rating: 'Rating',
        readNow: 'Read Now',
        addToShelf: 'Add to Shelf',
        totalChapters: 'Total {total} chapters',
        ascSort: 'Ascending',
        descSort: 'Descending',
        latest: 'Latest',
        noChapters: 'No chapters yet',
        authorOtherWorks: 'Other Works',
        similarRecommend: 'Similar Recommendations',
        more: 'More',
        books: 'books',
        bookExperience: 'Book Experience',
        reParse: 'Re-parse',
        reParseTitle: 'Re-parse Chapters',
        reParseConfirm: 'This will delete existing chapter index and re-parse from the content file. Continue?',
        reParseSuccess: 'Chapters re-parsed: {old} → {new}',
        reParseFailed: 'Failed to re-parse chapters',
      },
      reader: {
        prevChapter: 'Prev',
        nextChapter: 'Next',
        catalog: 'Catalog',
        detail: 'Detail',
        addShelf: 'Shelf',
        dayMode: 'Day',
        nightMode: 'Night',
        phone: 'Phone',
        words: 'Words',
        type: 'Type',
        status: 'Status',
        unknown: 'Unknown',
        ongoing: 'Ongoing',
        finished: 'Finished',
        writtenBy: 'by',
        loading: 'Loading...',
        noContent: 'No content'
      },
      catalog: {
        title: 'Catalog'
      }
    },
    admin: {
      system: {
        common: {
          status: {
            enable: 'Enable',
            disable: 'Disable'
          }
        },
        role: {
          title: 'Roles',
          roleName: 'Role Name',
          roleCode: 'Role Code',
          roleStatus: 'Status',
          roleDesc: 'Description',
          menuAuth: 'Menu Permissions',
          buttonAuth: 'Button Permissions',
          form: {
            roleName: 'Please enter role name',
            roleCode: 'Please enter role code',
            roleStatus: 'Please select status',
            roleDesc: 'Please enter description'
          },
          dataScope: {
            title: 'Data Scope',
            all: 'All',
            custom: 'Custom Dept',
            dept: 'Current Dept',
            deptAndSub: 'Dept & Sub',
            self: 'Self Only'
          },
          addRole: 'Add Role',
          editRole: 'Edit Role'
        },
        user: {
          title: 'Users',
          userName: 'User Name',
          userGender: 'Gender',
          nickName: 'Nick Name',
          userPhone: 'Phone',
          userEmail: 'Email',
          userStatus: 'Status',
          userRole: 'Role',
          form: {
            userName: 'Please enter user name',
            userGender: 'Please select gender',
            nickName: 'Please enter nick name',
            userPhone: 'Please enter phone',
            userEmail: 'Please enter email',
            userStatus: 'Please select status',
            userRole: 'Please select role'
          },
          addUser: 'Add User',
          editUser: 'Edit User',
          gender: {
            male: 'Male',
            female: 'Female'
          }
        },
        menu: {
          home: 'Home',
          title: 'Menus',
          id: 'ID',
          parentId: 'Parent ID',
          menuType: 'Menu Type',
          menuName: 'Menu Name',
          routeName: 'Route Name',
          routePath: 'Route Path',
          pathParam: 'Path Param',
          layout: 'Layout',
          page: 'Page Component',
          i18nKey: 'I18n Key',
          icon: 'Icon',
          localIcon: 'Local Icon',
          iconTypeTitle: 'Icon Type',
          order: 'Sort Order',
          constant: 'Constant Route',
          keepAlive: 'Keep Alive',
          href: 'External Link',
          hideInMenu: 'Hide In Menu',
          activeMenu: 'Active Menu',
          multiTab: 'Multi Tab',
          fixedIndexInTab: 'Fixed Tab Index',
          query: 'Route Query',
          button: 'Button',
          buttonCode: 'Button Code',
          buttonDesc: 'Button Description',
          menuStatus: 'Status',
          form: {
            home: 'Please select home',
            menuType: 'Please select menu type',
            menuName: 'Please enter menu name',
            routeName: 'Please enter route name',
            routePath: 'Please enter route path',
            pathParam: 'Please enter path param',
            layout: 'Please select layout',
            page: 'Please select page component',
            i18nKey: 'Please enter i18n key',
            icon: 'Please select icon',
            localIcon: 'Please select local icon',
            order: 'Please enter sort order',
            keepAlive: 'Please select keep alive',
            href: 'Please enter external link',
            hideInMenu: 'Please select hide in menu',
            activeMenu: 'Please enter active menu route name',
            multiTab: 'Please select multi tab',
            fixedInTab: 'Please select fixed in tab',
            fixedIndexInTab: 'Please enter fixed tab index',
            queryKey: 'Please enter query key',
            queryValue: 'Please enter query value',
            button: 'Please select button',
            buttonCode: 'Please enter button code',
            buttonDesc: 'Please enter button description',
            menuStatus: 'Please select status'
          },
          addMenu: 'Add Menu',
          editMenu: 'Edit Menu',
          addChildMenu: 'Add Child Menu',
          type: {
            directory: 'Directory',
            menu: 'Menu'
          },
          iconType: {
            iconify: 'Iconify Icon',
            local: 'Local Icon'
          }
        },
        dept: {
          title: 'Departments',
          deptName: 'Dept Name',
          deptCode: 'Dept Code',
          leader: 'Leader',
          parentId: 'Parent ID',
          sortOrder: 'Sort Order',
          status: 'Status',
          form: {
            parentId: 'Please select parent dept',
            deptName: 'Please enter dept name',
            deptCode: 'Please enter dept code',
            leader: 'Please enter leader',
            sortOrder: 'Please enter sort order',
            status: 'Please select status'
          },
          addDept: 'Add Dept',
          addChildDept: 'Add Child Dept',
          editDept: 'Edit Dept',
          deptId: "ID"
        },
        dict: {
          title: 'Dict',
          dictName: 'Dict Name',
          dictCode: 'Dict Code',
          dictDesc: 'Description',
          dictStatus: 'Status',
          form: {
            dictName: 'Please enter dict name',
            dictCode: 'Please enter dict code',
            dictDesc: 'Please enter description',
            dictStatus: 'Please select status'
          },
          addDict: 'Add Dict',
          editDict: 'Edit Dict'
        },
        dictItem: {
          title: 'Dict Items',
          itemLabel: 'Label',
          itemValue: 'Value',
          itemDesc: 'Description',
          sortOrder: 'Sort Order',
          itemStatus: 'Status',
          form: {
            itemLabel: 'Please enter label',
            itemValue: 'Please enter value',
            itemDesc: 'Please enter description',
            sortOrder: 'Please enter sort order',
            itemStatus: 'Please select status'
          },
          addDictItem: 'Add Dict Item',
          editDictItem: 'Edit Dict Item'
        },
        log: {
          title: 'Logs',
          loginLog: 'Login Log',
          operationLog: 'Operation Log',
          userName: 'User Name',
          loginIp: 'Login IP',
          loginType: 'Login Type',
          loginResult: 'Login Result',
          module: 'Module',
          action: 'Action',
          clientIp: 'Client IP',
          startTime: 'Start Time',
          endTime: 'End Time',
          form: {
            userName: 'Please enter user name',
            loginIp: 'Please enter IP',
            loginType: 'Please select login type',
            loginResult: 'Please select result',
            module: 'Please enter module',
            action: 'Please enter action',
            clientIp: 'Please enter client IP',
            startTime: 'Please enter start time',
            endTime: 'Please enter end time'
          }
        }
      },
      library: {
        bookCategory: {
          title: 'Book Categories',
          categoryName: 'Category Name',
          categoryCode: 'Category Code',
          description: 'Description',
          parentId: 'Parent ID',
          sortOrder: 'Sort Order',
          isHot: 'Hot Category',
          categoryStatus: 'Status',
          form: {
            categoryName: 'Please enter category name',
            categoryCode: 'Please enter category code',
            description: 'Please enter description',
            sortOrder: 'Please enter sort order',
            categoryStatus: 'Please select status'
          },
          addCategory: 'Add Category',
          addChildCategory: 'Add Child Category',
          editCategory: 'Edit Category'
        },
        bookTag: {
          title: 'Book Tags',
          tagName: 'Tag Name',
          description: 'Description',
          usageCount: 'Usage Count',
          form: {
            tagName: 'Please enter tag name',
            description: 'Please enter description'
          },
          addTag: 'Add Tag',
          editTag: 'Edit Tag'
        },
        book: {
          title: 'Book Management',
          bookName: 'Title',
          author: 'Author',
          cover: 'Cover',
          intro: 'Introduction',
          categoryId: 'Category',
          language: 'Language',
          serialStatus: 'Serial Status',
          visibility: 'Visibility',
          totalChapters: 'Chapters',
          totalWords: 'Words',
          listingStatus: 'Status',
          statusListed: 'Listed',
          statusUnlisted: 'Unlisted',
          statusReviewing: 'Reviewing',
          statusRejected: 'Rejected',
          avgRating: 'Rating',
          ratingCount: 'Votes',
          tags: 'Tags',
          form: {
            title: 'Please enter book title',
            author: 'Please enter author',
            cover: 'Please enter cover URL',
            intro: 'Please enter introduction',
            categoryId: 'Please select category',
            language: 'Please select language',
            serialStatus: 'Please select serial status',
            visibility: 'Please select visibility',
            tags: 'Please select tags',
          },
          addBook: 'Add Book',
          editBook: 'Edit Book',
          bookDetail: 'Book Detail',
          totalCategories: "",
          upload: "",
          scan: "",
          chapters: "",
          uploadFile: "",
          fileFormat: "",
          uploadSuccess: "",
          uploadTitle: "",
          uploadAuthor: "",
          confirmImport: "",
          importSuccess: "",
          scanSuccess: "",
          scanConfirm: "",
          chapterTitle: "",
          chapterNo: "",
          wordCount: "",
          chapterStatus: "",
          fileManage: "",
          uploadList: "",
          parseStatus: "",
          parsePending: "",
          parseProcessing: "",
          parseSuccess: "",
          parseFailed: "",
          selectFile: "",
          startUpload: "",
          reScan: "",
          scanResult: "",
          scanning: ""
        },
        bookChapterRule: {
          title: "",
          ruleName: "",
          scopeType: "",
          scopeGlobal: "",
          scopeBook: "",
          pattern: "",
          titleGroup: "",
          minChapterLen: "",
          maxChapterLen: "",
          priority: "",
          description: "",
          status: "",
          addRule: "",
          editRule: "",
          rulePreview: "",
          form: {
            ruleName: "",
            pattern: "",
            description: ""
          }
        },
        bookFilterRule: {
          title: "",
          ruleName: "",
          matchType: "",
          matchKeyword: "",
          matchRegex: "",
          pattern: "",
          action: "",
          actionReplace: "",
          actionBlock: "",
          actionMark: "",
          replacement: "",
          applyStage: "",
          stageInput: "",
          stageOutput: "",
          category: "",
          severity: "",
          severityLow: "",
          severityMedium: "",
          severityHigh: "",
          description: "",
          status: "",
          addRule: "",
          editRule: "",
          form: {
            ruleName: "",
            pattern: "",
            replacement: "",
            description: ""
          }
        }
      }
    },
  },
  form: {
    required: 'Cannot be empty',
    userName: {
      required: 'Please enter user name',
      invalid: 'User name format is incorrect'
    },
    phone: {
      required: 'Please enter phone number',
      invalid: 'Phone number format is incorrect'
    },
    pwd: {
      required: 'Please enter password',
      invalid: '6-18 characters, including letters, numbers, and underscores'
    },
    confirmPwd: {
      required: 'Please enter password again',
      invalid: 'The two passwords are inconsistent'
    },
    code: {
      required: 'Please enter verification code',
      invalid: 'Verification code format is incorrect'
    },
    email: {
      required: 'Please enter email',
      invalid: 'Email format is incorrect'
    }
  },
  dropdown: {
    closeCurrent: 'Close Current',
    closeOther: 'Close Other',
    closeLeft: 'Close Left',
    closeRight: 'Close Right',
    closeAll: 'Close All',
    pin: 'Pin Tab',
    unpin: 'Unpin Tab'
  },
  icon: {
    themeConfig: 'Theme Configuration',
    themeSchema: 'Theme Schema',
    lang: 'Switch Language',
    fullscreen: 'Fullscreen',
    fullscreenExit: 'Exit Fullscreen',
    reload: 'Reload Page',
    collapse: 'Collapse Menu',
    expand: 'Expand Menu',
    pin: 'Pin',
    unpin: 'Unpin'
  },
  datatable: {
    itemCount: 'Total {total} items',
    fixed: {
      left: 'Left Fixed',
      right: 'Right Fixed',
      unFixed: 'Unfixed'
    }
  }
};

export default local;
