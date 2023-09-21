import { useEffect, useState } from "react"
import url from "../main";

const Start = () => {

    const [data, setData] = useState<string>()
    const [pin, setPin] = useState<string>("0000");

    useEffect(() => {
        let subscribed = true;

        const getData = async() => {
            if (subscribed) {
                console.log("fetching data in Start.tsx ...")
                const response = await fetch(url+"/test");
                const json = await response.json();
                setData(json.data);
            }
        }
        
        getData();

        return () => {
            subscribed = false;
        }
    }, [])

    const handleForm = () => {
        console.log("submitting form...")
    }

    return (
        <div className="w-screen h-screen flex flex-col items-center justify-center p-4">
            {data}
            <form onSubmit={handleForm} className="max-w-[320px] max-h-[400px] bg-gray-800 rounded-lg flex items-center justify-center flex-col p-4">
                <label className="mb-2">Enter Game Pin:</label>
                <input className="bg-transparent px-2 py-1 text-center" type="text" value={pin} onChange={(e) => setPin(e.target.value)}/>
            </form>
        </div>
    )
}

export default Start