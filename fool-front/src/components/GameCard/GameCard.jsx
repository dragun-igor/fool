import React from 'react'
import './styles.css'

export const GameCard = (props) => {
    const { card, onClick } = props

    let className = "game-card"

    if (card.trump_suit) {
        className += " trump-suit"
    }

    if (card.selected) {
        className += " selected"
    }

    return (
        <div className={className} onClick={() => onClick(card.id)}>
            <h1>{card.denomination}</h1>
            <h1>{card.suit}</h1>
        </div>
    )
}