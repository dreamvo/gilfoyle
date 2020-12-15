<template>
  <v-app>
    <v-app-bar flat color="white" light>
      <v-img
        :src="require('./assets/logo.svg')"
        class="my-3"
        contain
        height="64"
      />

      <v-btn outlined color="#66f">
        <v-icon light>mdi-plus</v-icon>
        Create
      </v-btn>
    </v-app-bar>

    <v-app-bar flat color="#34495e" dark>
      <v-app-bar-nav-icon
        @click="drawerMenu = !drawerMenu"
      ></v-app-bar-nav-icon>

      <v-toolbar-title>Gilfoyle dashboard</v-toolbar-title>

      <v-spacer></v-spacer>
    </v-app-bar>

    <!--      <v-navigation-drawer-->
    <!--          v-model="drawerMenu"-->
    <!--          fixed-->
    <!--          clipped-->
    <!--      >-->
    <!--        <v-list-->
    <!--            nav-->
    <!--            dense-->
    <!--        >-->
    <!--          <v-list-item-group-->
    <!--              v-model="group"-->
    <!--              active-class="deep-purple&#45;&#45;text text&#45;&#45;accent-4"-->
    <!--          >-->
    <!--            <v-list-item>-->
    <!--              <v-list-item-icon>-->
    <!--                <v-icon>mdi-home</v-icon>-->
    <!--              </v-list-item-icon>-->
    <!--              <v-list-item-title>Home</v-list-item-title>-->
    <!--            </v-list-item>-->

    <!--            <v-list-item>-->
    <!--              <v-list-item-icon>-->
    <!--                <v-icon>mdi-account</v-icon>-->
    <!--              </v-list-item-icon>-->
    <!--              <v-list-item-title>Account</v-list-item-title>-->
    <!--            </v-list-item>-->
    <!--          </v-list-item-group>-->
    <!--        </v-list>-->
    <!--      </v-navigation-drawer>-->

    <v-main>
      <HelloWorld />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import HelloWorld from "./components/HelloWorld.vue";

export default Vue.extend({
  name: "App",
  components: {
    HelloWorld
  },
  data: () => ({
    drawerMenu: true
  })
});
</script>
