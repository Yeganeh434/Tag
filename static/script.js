
// register approved tag
document.getElementById('register_approved_tag_button').addEventListener('click', function() {
    document.getElementById('mainPage').classList.add('hidden');
    // document.getElementById('mainPage').style.display="none";
    document.getElementById('register_approved_tag').classList.remove('hidden');
});

document.getElementById('reg_a_form').addEventListener('submit', function(event) {
    event.preventDefault();

    const title = document.getElementById('reg_a_title').value;
    const description = document.getElementById('reg_a_description').value;
    const key=document.getElementById('reg_a_key').value;
    const picture=document.getElementById('reg_a_picture').value;


    fetch('/register_approved_tag', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ title:title , description:description , key:key, picture:picture})
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert(data.message)
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});


//register tag as draft
document.getElementById("register_tag_as_draft_button").addEventListener('click',function(){
    document.getElementById("mainPage").classList.add("hidden");
    document.getElementById("register_tag_as_draft").classList.remove("hidden");
})

document.getElementById("reg_d_form").addEventListener("submit",function(event){
    event.preventDefault();

    const title = document.getElementById('reg_d_title').value;
    const description = document.getElementById('reg_d_description').value;
    const key=document.getElementById('reg_d_key').value;
    const picture=document.getElementById('reg_d_picture').value;

    fetch("/register_tag_as_draft",{
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ title:title , description:description , key:key, picture:picture})
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert(data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})


//approve or reject tag
document.getElementById("approve_or_reject_tag_button").addEventListener('click',function(){
    document.getElementById("mainPage").classList.add("hidden");
    document.getElementById("approve_or_reject_tag").classList.remove("hidden");
})

document.getElementById("ap_or_rej_form").addEventListener("submit",function(event){
    event.preventDefault();

    const IdValue = document.getElementById('ap_or_rej_ID').value;
    let ID;
    ID=parseInt(IdValue);
    const approved = document.getElementById("approved_status");
    const rejected = document.getElementById("rejected_status");
    let status;
    if (approved.checked) {
        status=true;
    } else if (rejected.checked) {
        status=false;
    } else {
        alert("please select status of tag");
        return;
    }
   
    fetch("/approve_or_reject_tag",{
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id:ID , isApproved:status})
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert(data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})


//merge tags
document.getElementById("merge_tags_button").addEventListener('click',function(){
    document.getElementById("mainPage").classList.add("hidden");
    document.getElementById("merge_tags").classList.remove("hidden");
})

document.getElementById("merge_form").addEventListener("submit",function(event){
    event.preventDefault();

    const IdValue = document.getElementById('merge_ID').value;
    let ID;
    ID=parseInt(IdValue);
    const title = document.getElementById('merge_title').value;
    const description = document.getElementById('merge_description').value;
    const key=document.getElementById('merge_key').value;
    const picture=document.getElementById('merge_picture').value;

    fetch("/merge_tags",{
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ originalTagID:ID ,title:title , description:description , key:key, picture:picture})
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert(data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})


//register_tag_relationship
document.getElementById('register_tag_relationship_button').addEventListener('click', function() {
    document.getElementById('mainPage').classList.add('hidden');
    document.getElementById('register_tag_relationship').classList.remove('hidden');
});

document.getElementById('reg_rel_form').addEventListener('submit', function(event) {
    event.preventDefault();

    const fromTag = document.getElementById('reg_rel_fromTag').value;
    let fromTagID;
    fromTagID=parseInt(fromTag);
    const toTag = document.getElementById('reg_rel_toTag').value;
    let toTagID;
    toTagID=parseInt(toTag);
    const relationshipType=document.getElementById('reg_rel_type').value;
    
    fetch('/register_tag_relationship', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ fromTag:fromTagID , toTag:toTagID , relationshipType:relationshipType})
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert(data.message)
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});


//set_tag_relationship
document.getElementById("set_tag_relationship_button").addEventListener('click',function(){
    document.getElementById("mainPage").classList.add("hidden");
    document.getElementById("set_tag_relationship").classList.remove("hidden");
})

document.getElementById("set_rel_form").addEventListener("submit",function(event){
    event.preventDefault();

    const IdValue = document.getElementById('set_rel_ID').value;
    let ID;
    ID=parseInt(IdValue);
    const relationshipType = document.getElementById('set_rel_type').value;
   
    fetch("/set_tag_relationship",{
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id:ID , relationshipType:relationshipType})
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert( data.error);
        } else {
            alert(data.message);
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
})

//get_related_tags_by_key
document.getElementById('get_related_tags_by_key_button').addEventListener('click', function() {
    document.getElementById('mainPage').classList.add('hidden');
    document.getElementById('get_related_tags_by_key').classList.remove('hidden');
});

document.getElementById('by_key_form').addEventListener('submit', function(event) {
    event.preventDefault();

    const key=document.getElementById('by_key').value;

    fetch(`/get_related_tags_by_key/${key}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            document.getElementById('responseMessage_key').innerText = data.error;
        } else if(data.message) {
            document.getElementById('responseMessage_key').innerText = data.message;
        }else{
            document.getElementById('responseMessage_key').innerText = "IDs of tags related to this key:\n"+data; 
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});


//get_related_tags_by_ID
document.getElementById('get_related_tags_by_ID_button').addEventListener('click', function() {
    document.getElementById('mainPage').classList.add('hidden');
    document.getElementById('get_related_tags_by_ID').classList.remove('hidden');
});

document.getElementById('by_ID_form').addEventListener('submit', function(event) {
    event.preventDefault();

    const ID=document.getElementById('by_ID').value;

    fetch(`/get_related_tags_by_id/${ID}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            document.getElementById('responseMessage_ID').innerText = data.error;
        } else if(data.message) {
            document.getElementById('responseMessage_ID').innerText = data.message;
        }else{
            document.getElementById('responseMessage_ID').innerText = "IDs of tags related to this ID:\n"+data; 
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});



//get_related_tags_by_title
document.getElementById('get_related_tags_by_title_button').addEventListener('click', function() {
    document.getElementById('mainPage').classList.add('hidden');
    document.getElementById('get_related_tags_by_title').classList.remove('hidden');
});

document.getElementById('by_title_form').addEventListener('submit', function(event) {
    event.preventDefault();

    const title=document.getElementById('by_title').value;

    fetch(`/search_tag_by_title/${title}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            document.getElementById('responseMessage_title').innerText = data.error;
        } else if(data.message) {
            document.getElementById('responseMessage_title').innerText = data.message;
        }else{
            document.getElementById('responseMessage_title').innerText = "IDs of tags related to this title:\n"+data; 
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});