import axios from 'axios'

async function getPostsUsingFetch(){
    const response=await fetch('https://jsonplaceholder.typicode.com/posts')
    if (!response.ok) {
        throw new Error('Failed to fetch posts')
    }
    const data = await response.json()
    return data
}

async function getPostUsingAxios(){
    const res=await axios.get('https://jsonplaceholder.typicode.com/posts')
    return res.data
}

async function compareBothCalls(){
    const fetchPost=await getPostsUsingFetch()
    const axiosPost=await getPostUsingAxios()
    console.log("----------fetch result---------")
    console.log(fetchPost)
    console.log("----------axios result---------")
    console.log(axiosPost)  
}

compareBothCalls()