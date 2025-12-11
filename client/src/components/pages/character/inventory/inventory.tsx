/// TODO: 60 10lvl 120 20lvl 250 30lvl 500 40lvl 1000 50 lvl 2500 60lvl (3400 - 6200)
import {CurrencyDisplay} from "../../../tamplates/currencyDisplay/currencyDisplay.tsx";
import {InventorySlot} from "../../../tamplates/inventorySlot/inventorySlot.tsx";
import {mightySwordOfAzeroth} from "types/game/game-item/baseItem.ts";

const slotsInRow = 11;
const slotsInColumn = 3;

export const Inventory = () => {
    return (
        <div className="flex items-center justify-center w-full h-full ">
            <div className="aspect-[5/7] h-full bg-[#5F5F5F] flex flex-col gap-0.5 p-8 m-auto">
                <div className="flex basis-3/5 flex-grow-0 flex-row gap-2">
                    <div className="flex-[2] bg-black flex flex-col">
                        <div className="flex-[4%]">

                        </div>
                        <div className="flex-[69%] flex justify-between items-center p-1">
                            <div className="h-full flex flex-col gap-2 justify-center">
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-15 h-15"></div>
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-15 h-15"></div>
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-15 h-15"></div>
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-15 h-15"></div>
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2-500 w-15 h-15"></div>
                            </div>
                            <div className="h-full flex flex-col gap-2 justify-center">
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-15 h-15"></div>
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-15 h-15"></div>
                                <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-15 h-15"></div>
                            </div>
                        </div>
                        <div className="flex-[22%] flex justify-between items-center p-1">
                            <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-23 h-full">

                            </div>
                            <div className="bg-[#6B6B6B] border-[#7B7B7B] border-2 w-23 h-full">

                            </div>
                        </div>
                        <div className="flex-[5%] flex gap-1 p-1">
                            <CurrencyDisplay value={2}/>
                            <CurrencyDisplay value={2}/>
                            <CurrencyDisplay value={2}/>
                        </div>
                    </div>

                    <div className="flex-[1] flex flex-col gap-1">
                        <div className="bg-[#aaaaaa] flex-[27%] flex flex-col">
                            <div className="flex-[65%] flex items-center justify-center h-full">
                                <div
                                    className="relative bg-[#333] w-18 h-full flex flex-col items-center justify-start text-white font-light pt-3">
                                    <span className="text-lg leading-none co">Level</span>
                                    <span className="text-3xl leading-tight">40</span>

                                    <div className="absolute bottom-0 left-1/2 -translate-x-1/2
                    w-0 h-0
                    border-l-[43px] border-l-transparent
                    border-r-[43px] border-r-transparent
                    border-b-[22px] border-b-[#aaaaaa]">
                                    </div>
                                </div>
                            </div>
                            <div className="flex-[35%] flex flex-col items-center justify-start text-white font-light">
                                <span className="text-2xl leading-none text-[#353535]">Name</span>
                                <span className="text-lg leading-tight text-[#353535]">Mage</span>
                            </div>
                        </div>
                        <div className="flex-[67%] border-[#878787] bg-[#414141] border-1">
                        </div>
                        <div className="flex-[5%] border-[#878787] border-1">
                            <CurrencyDisplay value={32}/>
                        </div>
                    </div>
                </div>

                <div className="flex basis-2/5 flex-grow-0 bg-black">
                    <div
                        className="grid h-full w-full gap-0.5 p-0.5"
                        style={{
                            gridTemplateColumns: `repeat(${slotsInRow}, 1fr)`,
                            gridTemplateRows: `repeat(${slotsInColumn}, 1fr)`,
                            aspectRatio: "11/3"
                        }}
                    >
                        {Array.from({ length: slotsInRow * slotsInColumn }).map((_, i) =>
                            i === 0 ? (
                                <InventorySlot key={i} item={mightySwordOfAzeroth} />
                            ) : (
                                <InventorySlot key={i} item={undefined} />
                            )
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}