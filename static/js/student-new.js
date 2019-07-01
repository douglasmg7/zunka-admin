var app = new Vue({
  el: '#app',
  data: {
    test: 1
  },
  created() {
    console.log('created');
  },
  methods: {
    // Get products.
    saveStudent(){
      axios({
        method: 'post',
        url: `/student-save}`,
        headers:{'csrf-token' : csrfToken}
      })
      .then((res)=>{
        console.log(res.data.msg);
      })
      .catch((err)=>{
        console.error(`Error - saveStudent(), err: ${err}`);
      });
    },
  },   
});