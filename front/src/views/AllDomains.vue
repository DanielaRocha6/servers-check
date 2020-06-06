<template>
  <div>
    <div v-if="urls.length >0">
      <h2>You've already searched for:</h2>
      <b-list-group>
        <b-list-group-item v-for="item in urls" :key="item">
          <b-link :to="`/domains/${item}`" >
            http://{{item}}
          </b-link>
        </b-list-group-item>
      </b-list-group>
    </div>
    <div v-if="urls.length ===0">
      <h1>You haven't checked any domains yet!</h1>
      <router-link to="/">Check some in here!</router-link>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      urls: []
    }
  },
  methods: {
    getAll() {
      axios({ method: 'GET', url: 'http://localhost:8090/allDomains', headers: { 'content-type': 'text/plain' } })
        .then(result => {
          this.items = result.data.Items
          this.items.forEach(this.makeUrlSlice)
        })
        .catch(error => {
          /*eslint-disable*/
          console.error(error);
          /* eslint-enable */
        })
    },
    makeUrlSlice(item) {
      this.urls.push(item.Url)
    }
  },
  mounted () {
    this.getAll()
  }
}
</script>
