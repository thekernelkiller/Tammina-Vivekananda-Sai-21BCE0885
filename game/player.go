package game

type Player struct {
    ID        string
    Name      string
    Characters []*Character
    Placement []string
}

func (p *Player) GetCharacterByName(name string) *Character {
    for _, char := range p.Characters {
        if char.Name == name {
            return char
        }
    }
    return nil
}