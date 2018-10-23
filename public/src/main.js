import Vue from 'vue'

import App from './App.vue'
import { router } from './_helpers'
import { store } from './_store'

new Vue({ el: '#app', render: h => h(App), router, store });
