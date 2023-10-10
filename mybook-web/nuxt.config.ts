// https://nuxt.com/docs/api/configuration/nuxt-config

import { loadEnv } from 'vite'

const envScript = process.env.npm_lifecycle_script.split(' ')
const envName = envScript[envScript.length - 1] // 通过启动命令区分环境
const envData = loadEnv(envName, 'env') as unknown as VITE_ENV_CONFIG


export default defineNuxtConfig({
  // devtools: { enabled: true }
  runtimeConfig: {
    public: {
      baseUrl: envData  // env下读取的数据
    }
  },
  
   vite: {
     envDir: '~/env', // 指定env文件夹
   }

})
