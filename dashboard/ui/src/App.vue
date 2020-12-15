<template>
  <v-app>
    <v-toolbar flat color="white" light>
      <v-col cols="6" sm="4" md="2">
        <RouterLink to="/">
          <v-avatar tile height="32" width="auto">
            <img :src="require('@/assets/logo.svg')" height="48px" alt="logo" />
          </v-avatar>
        </RouterLink>
      </v-col>

      <v-spacer></v-spacer>

      <v-btn outlined color="#66f">
        <v-icon light>mdi-plus</v-icon>
        Create
      </v-btn>
    </v-toolbar>

    <v-toolbar flat color="#34495e" dark>
      <v-app-bar-nav-icon
        @click="drawerMenu = !drawerMenu"
      ></v-app-bar-nav-icon>

      <v-toolbar-title>Overview</v-toolbar-title>
    </v-toolbar>

    <v-navigation-drawer v-model="drawerMenu" absolute temporary>
      <v-list nav dense>
        <v-list-item-group active-class="deep-purple--text text--accent-4">
          <v-list-item
            v-for="(item, index) in navigation"
            @click="$router.push(item.link)"
            :key="index"
          >
            <v-list-item-icon>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-icon>
            <v-list-item-title>
              {{ item.title }}
            </v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-container>
        <v-row no-gutters>
          <v-col md="12" style="min-height:90vh;">
            <RouterView />
          </v-col>
        </v-row>
      </v-container>
    </v-main>

    <v-system-bar height="54px" dark color="#34495e" v-if="healthy">
      <v-icon size="8" color="green">mdi-circle</v-icon>
      <span>Instance status: Running</span>
      <v-spacer></v-spacer>
    </v-system-bar>
    <v-system-bar height="54px" dark color="#34495e" v-else>
      <v-icon size="8" color="red">mdi-circle</v-icon>
      <span>Instance status: Unavailable</span>
      <v-spacer></v-spacer>
    </v-system-bar>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "App",
  components: {},
  data: () => ({
    healthy: false,
    drawerMenu: false,
    navigation: [
      {
        title: "Medias",
        link: "/medias",
        icon: "mdi-home"
      },
      {
        title: "Metrics",
        link: "/metrics",
        icon: "mdi-home"
      },
      {
        title: "Settings",
        link: "/settings",
        icon: "mdi-home"
      }
    ]
  })
});
</script>
